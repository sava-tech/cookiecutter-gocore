package users

import (
	"github.com/gin-gonic/gin"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/handlers"
	// "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/auth/interfaces/http"
	// "github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/middleware"
)

// RegisterRoutes registers all auth routes with the gin engine
func RegisterRoutes(
	router *gin.RouterGroup,
	user *handlers.UserHandler,
	auth *handlers.AuthHandler,
	verification *handlers.VerificationHandler,
) {

	// Public routes (no authentication required)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/otp-email-login", auth.OtpEmailLogin)
		authGroup.POST("/verify-otp-email-login", auth.VerifyOtpEmailLogin)
		authGroup.POST("/verify-email", verification.VerifyEmail)
		authGroup.POST("/forgot-password", auth.ForgotPassword)
		authGroup.POST("/reset-password", auth.ResetPassword)
		authGroup.POST("/refresh-token", auth.RefreshToken)
	}

	// Protected routes (authentication required)
	accountGroup := router.Group("/account")
	{
		accountGroup.GET("/get-account", user.GetAccount)
		// accountGroup.PUT("/update-account", user.UpdateAccount)
		// accountGroup.PUT("/change-password", user.ChangePassword)
		// accountGroup.DELETE("/deactivate", user.DeactivateAccount)
	}
}

// RegisterPublicRoutes registers routes that don't need any auth at all
func RegisterPublicRoutes(router *gin.RouterGroup) {
	publicGroup := router.Group("/public")
	{
		publicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}
}
