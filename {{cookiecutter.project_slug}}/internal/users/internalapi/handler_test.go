package internalapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestRoutesRejectMissingOrWrongAPIKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	// The key check runs before any query, so no database is needed.
	RegisterRoutes(router, nil, "service-key")

	for _, path := range []string{"/internal/v1/users/count", "/internal/v1/users/registered"} {
		for _, key := range []string{"", "wrong-key"} {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			req.Header.Set("api-access-key", key)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, http.StatusUnauthorized, rec.Code, "%s with key %q", path, key)
		}
	}
}

func TestRoutesFailClosedWithoutConfiguredKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoutes(router, nil, "")

	req := httptest.NewRequest(http.MethodGet, "/internal/v1/users/count", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
