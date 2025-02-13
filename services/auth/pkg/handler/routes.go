package authhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/services/shared/middleware"
)

func RegisterRoutes(r *gin.RouterGroup, h *AuthHandler) {
	r.POST("/register", h.HandleCreateUser)
	r.POST("/login", h.HandleLogin)

	authenticated := r.Group("")
	authenticated.Use(middleware.AuthMiddleware())
	{
		authenticated.GET("/profile", h.HandleGetProfile)
		authenticated.POST("/logout", h.HandleLogout)
	}

}
