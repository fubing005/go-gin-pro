package services_api

import (
	"context"
	"fmt"
	"shalabing-gin/global"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type kafkaService struct{}

var KafkaService = new(kafkaService)

func (kafkaService *kafkaService) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (kafkaService *kafkaService) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (kafkaService *kafkaService) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message received: key=%s, value=%s", string(msg.Key), string(msg.Value))

		// 处理消息逻辑
		if err := processkafkaMessage(msg); err != nil {
			fmt.Printf("Failed to process message: %v", err)
			continue // 如果处理失败，不标记消息为已消费
		}

		/*
			sess.MarkMessage(msg, "") 的作用是告诉 Kafka：
				当前消费者已经成功处理了这条消息。
				将该消息的偏移量提交到消费者组的偏移量记录中（即 Committed Offset）
			手动提交的优点:
				确保消息在处理成功后才提交偏移量，避免因处理失败导致的数据丢失。
				提高对消息消费的控制，特别是对于 幂等性处理 或 分布式事务 场景。
			自动提交的对比:
				自动提交（config.Consumer.Offsets.AutoCommit.Enable = true）会在后台自动提交偏移量，但如果消息处理失败，偏移量已经提交，可能导致消息丢失。
		*/
		sess.MarkMessage(msg, "") //Sarama 框架中 手动提交消息偏移量 的方法，用于标记某条消息已经被成功消费
	}
	return nil
}

// Kafka consumer
func StartConsumer() {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新消息开始消费

	client, err := sarama.NewConsumerGroup(global.App.Config.Queue.Kafka.Brokers, global.App.Config.Queue.Kafka.GroupId, config)
	if err != nil {
		global.App.Log.Error("Failed to start Kafka consumer group: ", zap.Any("err", err))
	}
	defer client.Close()

	for {
		err := client.Consume(context.Background(), []string{global.App.Config.Queue.Kafka.Topic}, &kafkaService{})
		if err != nil {
			global.App.Log.Error("Error consuming messages: ", zap.Any("err", err))
		}
	}
}

// processMessage 模拟处理消息
func processkafkaMessage(msg *sarama.ConsumerMessage) error {
	global.App.Log.Info("Processing message: " + string(msg.Value))
	return nil
}
