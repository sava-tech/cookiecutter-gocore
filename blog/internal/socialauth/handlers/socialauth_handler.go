package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/michaelassa01/blog/internal/socialauth/application"
	"github.com/michaelassa01/blog/pkg/response"
)

type SocialAuthHandler struct {
	service *application.SocialAuthService
}

func NewSocialAuthHandler(service *application.SocialAuthService) *SocialAuthHandler {
	return &SocialAuthHandler{
		service: service,
	}
}

// Login initiates OAuth flow
// @Summary Login with social provider
// @Description Redirect to provider's OAuth login page
// @Tags Auth
// @Param provider path string true "Provider (google, apple)"
// @Router /social-auth/{provider}/login [get]
func (h *SocialAuthHandler) Login(c *gin.Context) {
	// Get provider from path param and add as query param
	provider := c.Param("provider")

	// Create a new request with the provider as query param
	c.Request.URL.RawQuery = "provider=" + provider
	
	// Let gothic handle it
	gothic.BeginAuthHandler(c.Writer, c.Request)

}

// Callback handles OAuth callback
// @Summary OAuth Callback
// @Description Handle OAuth callback and authenticate user
// @Tags Auth
// @Param provider path string true "Provider (google, github)"
// @Router /social-auth/{provider}/callback [get]
func (h *SocialAuthHandler) Callback(c *gin.Context) {
	provider := c.Param("provider")
	c.Request.URL.RawQuery = "provider=" + provider + "&" + c.Request.URL.RawQuery

	// Complete OAuth flow and get user info
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		fmt.Println("this goth user error ::", err)
		// If there's a specific error or the user denied access
		// TODO: proper error response
		c.Redirect(http.StatusFound, "/login?error=oauth_failed")
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	// Process user with your service
	authResp, err := h.service.ProcessSocialUser(
		c.Request.Context(),
		provider,
		gothUser,
		userAgent,
		clientIP,
	)
	if err != nil {
		fmt.Println("this is the goth authResp error", err)
		// Redirect to login with error
		// TODO: proper error response
		c.Redirect(http.StatusFound, "/login?error=social_auth_failed")
		return
	}

	// Return JSON response with tokens
	response.Success(c, "Social login successful", authResp)
}
