package comments

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *CommentHandler) {
	r.POST("", h.HandleCreateComment)
	r.GET("/:id", h.HandleGetComments)
	r.DELETE("/:id", h.HandleDeleteCategory)
}
