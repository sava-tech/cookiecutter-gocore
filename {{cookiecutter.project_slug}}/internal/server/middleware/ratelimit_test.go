package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func newRateLimitedRouter(t *testing.T, mw gin.HandlerFunc) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	require.NoError(t, router.SetTrustedProxies(nil))
	router.Use(mw)
	router.GET("/a", func(c *gin.Context) { c.Status(http.StatusOK) })
	router.GET("/b", func(c *gin.Context) { c.Status(http.StatusOK) })
	return router
}

func doRequest(router *gin.Engine, path, remoteAddr, forwardedFor string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.RemoteAddr = remoteAddr
	if forwardedFor != "" {
		req.Header.Set("X-Forwarded-For", forwardedFor)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestRateLimiterBlocksAfterBurst(t *testing.T) {
	router := newRateLimitedRouter(t, NewRateLimiter(rate.Every(time.Hour), 2).Middleware())

	require.Equal(t, http.StatusOK, doRequest(router, "/a", "1.1.1.1:1", "").Code)
	require.Equal(t, http.StatusOK, doRequest(router, "/a", "1.1.1.1:1", "").Code)

	rec := doRequest(router, "/a", "1.1.1.1:1", "")
	require.Equal(t, http.StatusTooManyRequests, rec.Code)
	require.NotEmpty(t, rec.Header().Get("Retry-After"))
	require.Contains(t, rec.Body.String(), "RATE_LIMIT_EXCEEDED")

	// A different client has its own budget.
	require.Equal(t, http.StatusOK, doRequest(router, "/a", "2.2.2.2:1", "").Code)
}

func TestRateLimiterIgnoresSpoofedForwardedFor(t *testing.T) {
	router := newRateLimitedRouter(t, NewRateLimiter(rate.Every(time.Hour), 1).Middleware())

	require.Equal(t, http.StatusOK, doRequest(router, "/a", "1.1.1.1:1", "9.9.9.1").Code)
	require.Equal(t, http.StatusTooManyRequests, doRequest(router, "/a", "1.1.1.1:1", "9.9.9.2").Code)
}

func TestPerRouteRateLimiterSeparatesRoutes(t *testing.T) {
	router := newRateLimitedRouter(t, NewRateLimiter(rate.Every(time.Hour), 1).PerRouteMiddleware())

	require.Equal(t, http.StatusOK, doRequest(router, "/a", "1.1.1.1:1", "").Code)
	require.Equal(t, http.StatusTooManyRequests, doRequest(router, "/a", "1.1.1.1:1", "").Code)
	require.Equal(t, http.StatusOK, doRequest(router, "/b", "1.1.1.1:1", "").Code)
}
