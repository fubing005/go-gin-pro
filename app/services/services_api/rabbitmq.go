package services_api

import (
	"errors"
	"log"
	"shalabing-gin/global"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type rabbitmqService struct{}

var RabbitmqService = new(rabbitmqService)

// PublishMessage 发布消息到 RabbitMQ
func (rabbitmqService *rabbitmqService) PublishMessage(exchange, routingKey, message string) (err error) {
	conn := global.App.Rabbitmq.GetConnection()
	defer global.App.Rabbitmq.ReleaseConnection(conn)

	ch, err := conn.Channel()
	if err != nil {
		err = errors.New("Failed to open a channel: " + err.Error())
		return
	}
	defer ch.Close()

	// 开启确认模式
	ch.Confirm(false)
	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation, 1))

	// 声明交换机
	err = ch.ExchangeDeclare(
		exchange,
		"direct", //direct,fanout,topic
		true,     // durable 持久化
		false,    // auto-deleted
		false,    // internal
		false,    // noWait
		nil,
	)
	if err != nil {
		err = errors.New("Failed to declare exchange:" + err.Error())
		return
	}

	// 发布消息
	err = ch.Publish(
		exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		err = errors.New("Failed to publish message:" + err.Error())
		return
	}

	// 等待确认
	select {
	case confirm := <-confirms:
		if confirm.Ack {
			global.App.Log.Info("Message confirmed by RabbitMQ")
			return nil
		}
		global.App.Log.Info("Message was not confirmed")
		return err
	case <-time.After(5 * time.Second):
		global.App.Log.Info("Message confirmation timeout")
		return err
	}
}

// StartConsumer 启动消费者
func (rabbitmqService *rabbitmqService) StartConsumer(queue, exchange, routingKey string) {
	for {
		conn := global.App.Rabbitmq.GetConnection()
		defer global.App.Rabbitmq.ReleaseConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.App.Log.Info("Failed to open a channel: " + err.Error())
			continue
		}
		defer ch.Close()

		// 声明队列
		_, err = ch.QueueDeclare(
			queue,
			true,  // durable
			false, // autoDelete
			false, // exclusive
			false, // noWait
			amqp091.Table{
				"x-dead-letter-exchange":    exchange + ".dlx",
				"x-dead-letter-routing-key": routingKey,
			},
		)
		if err != nil {
			global.App.Log.Info("Failed to declare queue: " + err.Error())
			continue
		}

		// 注册消费者
		msgs, err := ch.Consume(
			queue,
			"",
			false, // 手动 ACK
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			global.App.Log.Info("Failed to register a consumer: " + err.Error())
			continue
		}

		global.App.Log.Info("Consumer started...")
		for d := range msgs {
			err := processRabbitmqMessage(d.Body) // 模拟处理消息
			if err != nil {
				global.App.Log.Info("Failed to process message: " + err.Error())
				d.Nack(false, false) // 发送到死信队列
			} else {
				d.Ack(false) // 确认消费
			}
		}
		global.App.Log.Info("Consumer connection lost, reconnecting...")
	}
}

// processMessage 模拟处理消息
func processRabbitmqMessage(body []byte) error {
	log.Printf("Processing message: %s", body)
	return nil
}
