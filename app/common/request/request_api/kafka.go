package request_api

import (
	"shalabing-gin/app/common/request"
)

type KafkaRequest struct {
	Key   string `form:"key" json:"key" binding:"required"`
	Value string `form:"value" json:"value" binding:"required"`
}

func (kafkaRequest KafkaRequest) GetMessages() request.ValidatorMessages {
	return request.ValidatorMessages{
		"key.required":   "key is required",
		"value.required": "value is required",
	}
}
