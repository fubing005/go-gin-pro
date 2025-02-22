package services_api

import (
	"encoding/json"
	"fmt"
	"log"
	"shalabing-gin/app/models"
	"shalabing-gin/global"
)

type rabbitmqConsumerService struct{}

var RabbitmqConsumerService = new(rabbitmqConsumerService)

// 监听原始队列（订单状态队列）
func ConsumeMessages() {
	conn, ch, err := Connect()
	if err != nil {
		log.Fatal("RabbitMQ 连接失败:", err)
	}
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		"order_status_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("无法消费消息:", err)
	}

	for msg := range msgs {
		var data struct {
			OrderID uint               `json:"order_id"`
			Status  models.OrderStatus `json:"status"`
		}
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Println("消息解析失败:", err)
			continue
		}

		// 更新数据库状态
		var order models.Order
		if err := global.App.DB.First(&order, data.OrderID).Error; err == nil {
			order.Status = data.Status
			global.App.DB.Save(&order)
			fmt.Println("订单更新成功:", order)
		} else {
			log.Println("订单不存在:", data.OrderID)
		}
	}
}

// 监听死信队列（超时订单）
func ConsumeDeadLetterMessages() {
	conn, ch, err := Connect()
	if err != nil {
		log.Fatal("RabbitMQ 连接失败:", err)
	}
	defer conn.Close()
	defer ch.Close()

	// 消费死信队列
	msgs, err := ch.Consume(
		"order_status_dead_letter_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("无法消费死信消息:", err)
	}

	for msg := range msgs {
		var data struct {
			OrderID uint               `json:"order_id"`
			Status  models.OrderStatus `json:"status"`
		}
		if err := json.Unmarshal(msg.Body, &data); err != nil {
			log.Println("死信消息解析失败:", err)
			continue
		}

		// 处理超时订单（例如，将其状态更新为“已取消”）
		var order models.Order
		if err := global.App.DB.First(&order, data.OrderID).Error; err == nil {
			order.Status = models.Canceled
			global.App.DB.Save(&order)
			fmt.Println("死信队列处理：订单超时取消:", order)
		} else {
			log.Println("死信队列处理失败：订单不存在:", data.OrderID)
		}
	}
}
