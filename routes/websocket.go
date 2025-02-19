package routes

import (
	"shalabing-gin/app/controllers/websocket"

	"github.com/gin-gonic/gin"
)

func SetWebsocketGroupRoutes(router *gin.RouterGroup) {
	router.GET("/chat", websocket.ChatController{}.ChatMessage)
}
