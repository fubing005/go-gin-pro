package websocket

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/common/request/request_websocket"
	"shalabing-gin/app/common/response"
	"shalabing-gin/app/services/services_websocket"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
}

// ChatMessage 处理 WebSocket 连接
func (con ChatController) ChatMessage(c *gin.Context) {
	var form request_websocket.ChatRequest
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	services_websocket.HandleChatMessage(c, form)
}
