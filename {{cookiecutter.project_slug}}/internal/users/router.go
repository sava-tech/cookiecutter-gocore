package users

import "github.com/gin-gonic/gin"


func RegisterUsersRoutes(r *gin.RouterGroup, h *Handler) {
    group := r.Group("/users")
    {
        group.POST("/users", h.CreateUser)
        // group.GET("/:id", h.getPost)
        // group.PUT("/:id", h.updatePost)
        // group.DELETE("/:id", h.deletePost)
    }
}


