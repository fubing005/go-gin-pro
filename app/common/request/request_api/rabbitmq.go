package request_api

import (
	"shalabing-gin/app/common/request"
)

type RabbitMQRequest struct {
	Exchange   string `form:"exchange" json:"exchange" binding:"required"`
	RoutingKey string `form:"routing_key" json:"routing_key" binding:"required"`
	Message    string `form:"message" json:"message" binding:"required"`
}

func (rabbitMQRequest RabbitMQRequest) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"exchange.required":    "exchange is required",
		"routing_key.required": "routing_key is required",
		"message.required":     "message is required",
	}
}
