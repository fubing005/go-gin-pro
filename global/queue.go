package global

import (
	"sync"

	"github.com/IBM/sarama"
)

// QueuePool RabbitMQ 连接池
type Queue struct {
	connPool sync.Pool
	kafka    sarama.SyncProducer
}
