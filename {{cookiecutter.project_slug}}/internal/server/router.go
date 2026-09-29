package server

import (
	"log"
	"strings"

	// "github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

	// import modules router
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/server/middleware"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth"
	socialauthHandler "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/handlers"
	socialauthInfra "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/socialauth/infrastructure"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/handlers"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/utils"
)

func (s *Server) setupRouter(config utils.Config) {
	// IMPORTANT: Setup session store FIRST
	socialauthInfra.SetupSession(config.SessionSecret)

	// Then setup Goth providers
	gothConfig := socialauthInfra.GothProviderConfig{
		GoogleKey:    config.GoogleKey,
		GoogleSecret: config.GoogleSecret,
		AppleKey:     config.AppleKey,
		AppleSecret:  config.AppleSecret,
		CallbackURL:  config.SocialCallbackURL,
	}
	socialauthInfra.SetupGoth(gothConfig)

	router := gin.Default()
	// Only trust X-Forwarded-For from our own proxies; otherwise any client
	// could spoof c.ClientIP() and dodge per-IP rate limits.
	if err := router.SetTrustedProxies(splitCSV(config.TrustedProxies)); err != nil {
		log.Fatalf("invalid TRUSTED_PROXIES: %v", err)
	}
	router.Use(newGlobalRateLimiter(config).Middleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Stricter per-endpoint limit for credential/OTP routes (brute-force targets).
	authRateLimit := newAuthRateLimiter(config).PerRouteMiddleware()

	// Group that ONLY requires API Key
	apiKeyOnly := router.Group("/api/v1")
	// apiKeyOnly.Use(middleware.APIKeyAuth(config.ApiAccessKey))
	publicGroup := router.Group("/public")
	{
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
	// User routes
	usersHandler := handlers.NewUserHandler(s.UserServices.User)
	authHandler := handlers.NewAuthHandler(s.UserServices.Auth)
	verificationHandler := handlers.NewVerificationHandler(s.UserServices.Auth)

	// socialAuth routes
	socialHandler := socialauthHandler.NewSocialAuthHandler(s.SocialAuthServices.SocialAuth)

	// Registered Routes
	users.RegisterRoutes(
		apiKeyOnly,
		usersHandler,
		authHandler,
		verificationHandler,
		authRateLimit,
	)

	socialauth.RegisterRoutes(
		apiKeyOnly,
		socialHandler,
	)

	s.router = router
}

// newGlobalRateLimiter limits every request per client IP
// (default 5 req/s, burst 10).
func newGlobalRateLimiter(config utils.Config) *middleware.RateLimiter {
	rps, burst := config.RateLimitRPS, config.RateLimitBurst
	if rps <= 0 {
		rps = 5
	}
	if burst <= 0 {
		burst = 10
	}
	return middleware.NewRateLimiter(rate.Limit(rps), burst)
}

// newAuthRateLimiter limits login/OTP/password endpoints per client IP and
// endpoint (default 5 req/min, burst 5).
func newAuthRateLimiter(config utils.Config) *middleware.RateLimiter {
	perMinute, burst := config.AuthRateLimitPerMinute, config.AuthRateLimitBurst
	if perMinute <= 0 {
		perMinute = 5
	}
	if burst <= 0 {
		burst = 5
	}
	return middleware.NewRateLimiter(rate.Limit(perMinute/60), burst)
}

// splitCSV splits a comma-separated list, trimming blanks. Returns nil for
// an empty string.
func splitCSV(value string) []string {
	var items []string
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}
