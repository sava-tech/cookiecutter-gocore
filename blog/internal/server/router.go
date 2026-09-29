package server

import (

	// "github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// import modules router
	"github.com/michaelassa01/blog/internal/server/middleware"
	"github.com/michaelassa01/blog/internal/socialauth"
	socialauthHandler "github.com/michaelassa01/blog/internal/socialauth/handlers"
	socialauthInfra "github.com/michaelassa01/blog/internal/socialauth/infrastructure"
	"github.com/michaelassa01/blog/internal/users"
	"github.com/michaelassa01/blog/internal/users/handlers"
	"github.com/michaelassa01/blog/utils"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.RateLimitMiddleware())

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
	)

	socialauth.RegisterRoutes(
		apiKeyOnly,
		socialHandler,
	)

	s.router = router
}
