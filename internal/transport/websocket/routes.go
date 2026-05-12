package websocket

import "github.com/gin-gonic/gin"

// RegisterRoutes registers WebSocket notification routes.
func RegisterRoutes(router *gin.Engine, hub *Hub) {
	if hub == nil {
		return
	}

	router.GET("/ws/notifications", hub.HandleNotifications)
}
