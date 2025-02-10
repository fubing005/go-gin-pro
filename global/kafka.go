package global

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

// Initialize Kafka producer
func InitProducer() *Queue {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(App.Config.Queue.Kafka.Brokers, config)
	if err != nil {
		App.Log.Error("Failed to create kafka producer: ", zap.Any("err", err))
	} else {
		App.Log.Error("success to create kafka producer: ", zap.Any("err", err))
	}
	return &Queue{
		kafka: producer,
	}
}

func (qp *Queue) GetProducer() sarama.SyncProducer {
	return qp.kafka
}
