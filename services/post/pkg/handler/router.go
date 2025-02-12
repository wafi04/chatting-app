package posthandler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *PostHandler) {
	r.POST("", h.HandleCreatePost)
	r.GET("/all", h.HandleGetAllPosts)
	r.GET("/user", h.HandleGetPostByUser)
}
