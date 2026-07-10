package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/application"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/internal/users/dto"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/response"
)

// UserHandler binds Gin routes to the service.
type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(s *application.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// GetAccount godoc
// @Summary Get user account
// @Description Retrieve user account information
// @Tags Account
// @Accept json
// @Produce json
// @Param request body dto.GetUserReq true "Get user account details"
// @Success 200 {object} response.Response{data=dto.GetUserRes} "Account retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Router /account/get-account [get]
func (h *UserHandler) GetAccount(c *gin.Context) {
	var req dto.GetUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request", err.Error())
		return
	}

	resp, err := h.service.GetAccount(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, "Failed to get account")
		return
	}
	response.Success(c, "Account retrieved successfully", resp)
}
