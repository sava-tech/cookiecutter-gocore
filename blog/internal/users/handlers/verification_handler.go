package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/michaelassa01/blog/internal/users/application"
	"github.com/michaelassa01/blog/internal/users/domain"
	"github.com/michaelassa01/blog/internal/users/dto"
	"github.com/michaelassa01/blog/pkg/response"
)

// VerificationHandler binds Gin routes to the service.
type VerificationHandler struct {
	service *application.AuthService
}

func NewVerificationHandler(s *application.AuthService) *VerificationHandler {
	return &VerificationHandler{service: s}
}

// VerifyEmail godoc
// @Summary Verify user email
// @Description Verify user email with verification code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyEmailRequest true "Verification details"
// @Success 200 {object} response.Response "Email verified successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 404 {object} response.Response "Invalid or expired verification code"
// @Router /auth/verify-email [post]
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	if err := h.service.VerifyEmail(c.Request.Context(), &req); err != nil {
		switch err {
		case domain.ErrInvalidVerificationCode:
			response.BadRequest(c, "Invalid or expired verification code", nil)
		default:
			response.InternalError(c, "Failed to verify email")
		}
		return
	}

	response.Success(c, "Email verified successfully", nil)
}
