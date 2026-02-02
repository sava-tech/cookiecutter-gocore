package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"{{ cookiecutter.module_path }}/pkg/token"
	u "{{ cookiecutter.module_path }}/utils"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// Each user/device gets its own limiter
var limiters = make(map[string]*rate.Limiter)
var mu sync.Mutex

// NewLimiter creates a limiter for each key (e.g. IP, DeviceID)
func NewLimiter(rps float64, burst int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(rps), burst)
}

func getLimiter(key string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[key]
	if !exists {
		// Example: 5 requests per second, burst up to 10
		limiter = NewLimiter(5, 10)
		limiters[key] = limiter
	}
	return limiter
}

// RateLimitMiddleware applies per-user or per-IP rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// You can also use c.ClientIP() or a user token instead of IP
		clientID := c.ClientIP()

		limiter := getLimiter(clientID)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, slow down.",
			})
			return
		}

		c.Next()
	}
}

func AuthMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeaderKey := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeaderKey) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, u.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeaderKey)
		if len(fields) < 2 {
			err := errors.New("invalide authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, u.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, u.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, u.ErrorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()

	}
}

// RoleMiddleware this implements the roles permissions
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		val, exists := ctx.Get("authorization_payload")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing payload"})
			return
		}

		payload, ok := val.(*token.Payload)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid payload type"})
			return
		}

		userRole := payload.AccountType

		// Allow all roles
		for _, r := range allowedRoles {
			if r == "any" {
				ctx.Next()
				return
			}
		}

		// Check matches
		for _, r := range allowedRoles {
			if r == userRole {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
	}
}

func APIKeyAuth(expectedKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("api-access-key")

		if apiKey == "" || apiKey != expectedKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - missing or invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
