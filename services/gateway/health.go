package gateway

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wafi04/chatting-app/services/shared/pkg/response"
)

func CheckCoon(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		healthStatus := true
		timestamp := time.Now().Format(time.RFC3339)
		data := struct {
			Health    bool   `json:"health"`
			Timestamp string `json:"time"`
		}{
			Health:    healthStatus,
			Timestamp: timestamp,
		}

		if healthStatus {
			response.SendSuccessResponse(c, http.StatusOK, "Connection Ready ", data)
		} else {
			response.SendErrorResponse(c, http.StatusServiceUnavailable, "Service Unhealthy")
		}
	})

}
