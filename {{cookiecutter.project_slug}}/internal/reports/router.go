// Package reports is a standalone service that serves user reports. It holds
// no user data: it verifies the caller's access token locally and fetches the
// data from the users service over HTTP.
package reports

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"{{ cookiecutter.module_path }}/internal/server/middleware"
	"{{ cookiecutter.module_path }}/pkg/response"
	"{{ cookiecutter.module_path }}/pkg/token"
	"{{ cookiecutter.module_path }}/utils"
)

const defaultLimit = 20

// UsersService is the data this service needs from the users service
type UsersService interface {
	CountUsers(ctx context.Context) (int64, error)
	ListRegisteredUsers(ctx context.Context, limit, offset int) ([]RegisteredUser, error)
}

type handler struct {
	users UsersService
}

type listRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

// NewRouter creates the reports HTTP router. Any signed-in user can read the
// count; the user list contains personal data, so it requires an admin.
func NewRouter(tokenMaker token.Maker, users UsersService) *gin.Engine {
	h := &handler{users: users}

	router := gin.Default()
	router.GET("/public/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	group := router.Group("/api/v1/reports/users", middleware.AuthMiddleWare(tokenMaker))
	{
		group.GET("/count", h.countUsers)
		group.GET("/registered", middleware.RoleMiddleware("admin"), h.listRegisteredUsers)
	}
	return router
}

func (h *handler) countUsers(c *gin.Context) {
	count, err := h.users.CountUsers(c.Request.Context())
	if err != nil {
		upstreamError(c, err)
		return
	}
	response.Success(c, "Users count retrieved", gin.H{"count": count})
}

func (h *handler) listRegisteredUsers(c *gin.Context) {
	var req listRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	if req.Limit == 0 {
		req.Limit = defaultLimit
	}

	registered, err := h.users.ListRegisteredUsers(c.Request.Context(), req.Limit, req.Offset)
	if err != nil {
		upstreamError(c, err)
		return
	}
	response.Success(c, "Registered users retrieved", gin.H{"users": registered})
}

func upstreamError(c *gin.Context, err error) {
	log.Printf("users service call failed: %v", err)
	c.JSON(http.StatusBadGateway, utils.ErrorResponseMessage("users service unavailable"))
}
