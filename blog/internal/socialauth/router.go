package socialauth

import (
	"github.com/gin-gonic/gin"
	"github.com/michaelassa01/blog/internal/socialauth/domain"
	"github.com/michaelassa01/blog/internal/socialauth/handlers"
	"github.com/michaelassa01/blog/pkg/response"
)

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *handlers.SocialAuthHandler,
) {
	socialAuthGroup := router.Group("/social-auth")
	socialAuthGroup.Use(validateProvider)
	{
		socialAuthGroup.GET("/:provider/login", handler.Login)
		socialAuthGroup.GET("/:provider/callback", handler.Callback)
	}
}

// validateProvider middleware validates the OAuth provider
func validateProvider(c *gin.Context) {
	provider := c.Param("provider")
	if !domain.IsValidProvider(provider) {
		response.BadRequest(c, "Invalid provider", "Supported providers: google, apple")
		c.Abort()
		return
	}
	c.Next()
}