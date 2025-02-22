package services_api

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type rabbitmqProducerService struct{}

var RabbitmqProducerService = new(rabbitmqProducerService)

// 连接 RabbitMQ
func Connect() (*amqp091.Connection, *amqp091.Channel, error) {
	conn, err := amqp091.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// 发布订单状态更新消息
func (rabbitmqProducerService *rabbitmqProducerService) PublishOrderUpdate(orderID uint, status string) (err error) {
	conn, ch, err := Connect()
	if err != nil {
		log.Println("RabbitMQ 连接失败:", err)
		return
	}
	defer conn.Close()
	defer ch.Close()

	// 声明订单状态队列，设置死信队列的路由键
	_, err = ch.QueueDeclare(
		"order_status_queue",
		true,  // 持久化
		false, // 自动删除
		false, // 排他
		false, // 非阻塞
		amqp091.Table{
			"x-dead-letter-exchange":    "", // 默认交换机
			"x-dead-letter-routing-key": "order_timeout",
		},
	)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(`{"order_id":%d, "status":"%s"}`, orderID, status)
	err = ch.Publish(
		"",
		"order_status_queue",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		log.Println("消息发送失败:", err)
	}

	return nil
}

// 发送超时订单消息到死信队列
func SendToDeadLetterQueue(message string) error {
	conn, ch, err := Connect()
	if err != nil {
		log.Println("RabbitMQ 连接失败:", err)
		return err
	}
	defer conn.Close()
	defer ch.Close()

	// 创建死信队列（DLX）
	dlQueueName := "order_status_dead_letter_queue"
	_, err = ch.QueueDeclare(
		dlQueueName,
		true,  // 是否持久化
		false, // 是否自动删除
		false, // 是否排他
		false, // 是否非阻塞
		amqp091.Table{
			"x-dead-letter-exchange":    "",              // 默认交换机
			"x-dead-letter-routing-key": "order_timeout", // 死信路由键
		},
	)
	if err != nil {
		return err
	}

	// 将超时订单消息发送到死信队列
	err = ch.Publish(
		"",
		"order_status_dead_letter_queue", // 死信队列名
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
		},
	)
	if err != nil {
		log.Println("无法发送消息到死信队列:", err)
		return err
	}

	log.Printf("消息成功发送到死信队列: %s\n", message)
	return nil
}
