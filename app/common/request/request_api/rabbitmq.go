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

type RabbitMQRequestOrderCreate struct {
	Amount    float64            `form:"amount" json:"amount"  binding:"required"`
	Status    models.OrderStatus `form:"status" json:"status"  binding:"required,oneof=pending paid shipped completed canceled"`
	UserID    uint               `form:"user_id" json:"user_id"  binding:"required"`
	ProductID uint               `form:"product_id" json:"product_id"  binding:"required"`
}

func (rabbitMQRequestOrderCreate RabbitMQRequestOrderCreate) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"amount.required":     "amount is required",
		"status.required":     "status is required",
		"status.oneof":        "status must be one of pending paid shipped completed canceled",
		"user_id.required":    "user_id is required",
		"product_id.required": "product_id is required",
	}
}

type RabbitMQRequestOrderStatusUpdate struct {
	ID        uint               `form:"id" json:"id" binding:"omitempty"` // 绑定 URL 查询参数 id
	NewStatus models.OrderStatus `form:"new_status" json:"new_status"  binding:"omitempty"`
}

func (rabbitMQRequestOrderStatusUpdate RabbitMQRequestOrderStatusUpdate) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"id.required":         "id is required",
		"new_status.required": "new_status is required",
	}
}
