package handlers

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/application"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/domain"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/dto"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/response"
)

// AuthHandler binds Gin routes to the service.
type AuthHandler struct {
	service *application.AuthService
}

func NewAuthHandler(s *application.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=dto.RegisterResponse} "Account created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 409 {object} response.Response "Email or phone already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrEmailAlreadyExists:
			response.Conflict(c, "Email already exists")
		case domain.ErrPhoneAlreadyExists:
			response.Conflict(c, "Phone number already exists")
		default:
			response.InternalError(c, "Failed to create account")
		}
		return
	}

	response.Created(c, "Account created successfully", user)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user and return access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login details"
// @Success 200 {object} response.Response{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Invalid email or password"
// @Failure 403 {object} response.Response "Email not verified"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	resp, err := h.service.Login(c.Request.Context(), &req, userAgent, clientIP)
	if err != nil {
		fmt.Println("server is not working", err)
		switch err {
		case domain.ErrInvalidCredentials:
			response.Unauthorized(c, "Invalid email or password")
		case domain.ErrEmailNotVerified:
			response.Forbidden(c, "Email not verified. Please check your email")
		default:
			response.InternalError(c, "Failed to login")
		}
		return
	}

	response.Success(c, "Login successful", resp)
}

// OtpEmailLogin godoc
// @Summary Otp Email Login a user
// @Description Authenticate user using the registered email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.OtpEmailLoginReq true "Otp Email Login"
// @Success 200
// @Router /auth/otp-email-login [post]
func (h *AuthHandler) OtpEmailLogin(c *gin.Context) {
	var req dto.OtpEmailLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	err := h.service.OtpEmailLogin(c.Request.Context(), &req)
	if err != nil {
		return
	}

	response.Success(c, "If your email is registered, you will receive a verification code", nil)
}

// VerifyOtpEmailLogin godoc
// @Summary verify Otp Email Login a user
// @Description Authenticate user using the registered email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOtpEmailLoginReq true "Verify Otp Email Login"
// @Success 200
// @Router /auth/verify-otp-email-login [post]
func (h *AuthHandler) VerifyOtpEmailLogin(c *gin.Context) {
	var req dto.VerifyOtpEmailLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}
	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	resp, err := h.service.VerifyOtpEmailLogin(c.Request.Context(), &req, userAgent, clientIP)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			response.Unauthorized(c, "Invalid email or code")
		case domain.ErrEmailNotVerified:
			response.Forbidden(c, "Email not verified. Please check your email")
		default:
			response.InternalError(c, "Failed to login")
		}
		return
	}

	response.Success(c, "Login successful", resp)

}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Initiate forgot password process
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "Forgot password details"
// @Success 200 {object} response.Response "If your email is registered, you will receive a verification code"
// @Failure 400 {object} response.Response "Invalid request"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	if err := h.service.ForgotPassword(c.Request.Context(), &req); err != nil {
		response.InternalError(c, "Failed to process request")
		return
	}

	response.Success(c, "If your email is registered, you will receive a verification code", nil)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password with verification code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset password details"
// @Success 200 {object} response.Response "Password reset successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Invalid or expired verification code"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), &req); err != nil {
		response.InternalError(c, "Failed to reset password")
		return
	}

	response.Success(c, "Password reset successfully", nil)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token details"
// @Success 200 {object} response.Response "Token refreshed successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 401 {object} response.Response "Invalid or expired refresh token"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}
	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		fmt.Println("Error refreshing token:", err)
		response.InternalError(c, "Failed to refresh token")
		return
	}
	response.Success(c, "Token refreshed successfully", resp)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user by invalidating the session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Logged out successfully"
// @Failure 401 {object} response.Response "Unauthorized - No session found"
// @Failure 401 {object} response.Response "Invalid session ID format"
// @Failure 500 {object} response.Response "Failed to logout"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get the session ID from the context (set by authentication middleware)
	// Optional: Get other user info for logging or validation
	// userEmail, _ := c.Get("user_email")
	// accountType, _ := c.Get("account_type")

	sessionIDValue, exists := c.Get("session_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized - No session found")
		return
	}

	// Convert the session ID to UUID
	sessionID, ok := sessionIDValue.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid session ID format")
		return
	}

	// Call service to delete the session
	if err := h.service.Logout(c.Request.Context(), sessionID); err != nil {
		// Check if session not found
		switch {
		case errors.Is(err, domain.ErrSessionNotFound):
			response.Unauthorized(c, "Unauthorized - Session not found")
		default:
			response.InternalError(c, "Failed to logout")
		}
		// Other errors
		response.InternalError(c, "Failed to logout")
		return
	}

	// Return success response
	response.Success(c, "Logged out successfully", nil)
}

// Account related handlers
