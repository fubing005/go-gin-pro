package request_websocket

import (
	"shalabing-gin/app/common/request"
)

type ChatRequest struct {
	ID string `form:"id" json:"id" binding:"required"`
}

func (chatRequest ChatRequest) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required": "用户 ID 不能为空",
	}
}
