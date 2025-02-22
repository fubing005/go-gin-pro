package request_api

import (
	"shalabing-gin/app/common/request"
	"shalabing-gin/app/models"
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

type RabbitMQRequestOrder struct {
	Amount    float64            `form:"amount" json:"amount"  binding:"omitempty"`
	Status    models.OrderStatus `form:"status" json:"status"  binding:"omitempty"`
	NewStatus models.OrderStatus `form:"new_status" json:"new_status"  binding:"omitempty"`
}

func (rabbitMQRequestOrder RabbitMQRequestOrder) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		// "id.required": "id is required",
	}
}

type RabbitMQRequestOrderID struct {
	ID uint `form:"id" json:"id" binding:"omitempty"` // 绑定 URL 查询参数 id
}

func (rabbitMQRequestOrderID RabbitMQRequestOrderID) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required": "id is required",
	}
}
