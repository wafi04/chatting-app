package authhandler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *AuthHandler) {
	r.POST("/register", h.HandleCreateUser)
	r.POST("/login", h.HandleLogin)

	r.GET("/profile", h.HandleGetProfile)
	r.POST("/logout", h.HandleLogout)

}
