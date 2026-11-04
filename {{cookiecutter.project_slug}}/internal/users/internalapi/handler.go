// Package internalapi exposes user data to other services. Its routes are
// authenticated with the shared API_ACCESS_KEY, not with end-user tokens.
package internalapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"{{ cookiecutter.module_path }}/internal/server/middleware"
	"{{ cookiecutter.module_path }}/internal/users/models"
	"{{ cookiecutter.module_path }}/utils"
)

const defaultLimit = 20

type handler struct {
	queries *users.Queries
}

type listRequest struct {
	Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int32 `form:"offset" binding:"omitempty,min=0"`
}

// RegisterRoutes registers the service-to-service routes
func RegisterRoutes(router gin.IRouter, db *pgxpool.Pool, apiKey string) {
	h := &handler{queries: users.New(db)}

	group := router.Group("/internal/v1", middleware.APIKeyAuth(apiKey))
	{
		group.GET("/users/count", h.countUsers)
		group.GET("/users/registered", h.listRegisteredUsers)
	}
}

func (h *handler) countUsers(c *gin.Context) {
	count, err := h.queries.CountUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseMessage("could not count users"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
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

	registered, err := h.queries.ListRegisteredUsers(c.Request.Context(), users.ListRegisteredUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponseMessage("could not list users"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": registered})
}
