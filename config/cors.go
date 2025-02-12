package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUpCors(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://192.168.100.9:3000")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	})
}
