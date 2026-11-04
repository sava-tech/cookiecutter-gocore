package reports

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"{{ cookiecutter.module_path }}/pkg/token"
	"{{ cookiecutter.module_path }}/utils"
)

const testAPIKey = "service-key"

// newUsersServer fakes the users service and enforces the service key
func newUsersServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("api-access-key") != testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/internal/v1/users/count":
			json.NewEncoder(w).Encode(map[string]int64{"count": 3})
		case "/internal/v1/users/registered":
			user := RegisteredUser{Email: "a@example.com"}
			json.NewEncoder(w).Encode(map[string][]RegisteredUser{"users": {user}})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestReportsAuthAndData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	maker, err := token.NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	tokenFor := func(accountType string) string {
		signed, _, err := maker.CreateToken("a@example.com", accountType, time.Minute)
		require.NoError(t, err)
		return signed
	}

	usersServer := newUsersServer(t)
	wrongKeyClient := NewUsersClient(usersServer.URL, "wrong-key")
	client := NewUsersClient(usersServer.URL, testAPIKey)

	testCases := []struct {
		name   string
		path   string
		token  string
		client *UsersClient
		code   int
		body   string
	}{
		{"count without token", "/api/v1/reports/users/count", "", client, http.StatusUnauthorized, ""},
		{"count with garbage token", "/api/v1/reports/users/count", "garbage", client, http.StatusUnauthorized, ""},
		{"count as user", "/api/v1/reports/users/count", tokenFor("user"), client, http.StatusOK, `"count":3`},
		{"list as user", "/api/v1/reports/users/registered", tokenFor("user"), client, http.StatusForbidden, ""},
		{"list as admin", "/api/v1/reports/users/registered", tokenFor("admin"), client, http.StatusOK, "a@example.com"},
		{"list with bad limit", "/api/v1/reports/users/registered?limit=1000", tokenFor("admin"), client, http.StatusBadRequest, ""},
		{"users service rejects key", "/api/v1/reports/users/count", tokenFor("user"), wrongKeyClient, http.StatusBadGateway, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := NewRouter(maker, tc.client)
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tc.code, rec.Code, rec.Body.String())
			require.Contains(t, rec.Body.String(), tc.body)
		})
	}
}
