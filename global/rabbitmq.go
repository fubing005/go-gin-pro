package global

import (
	"strconv"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// 初始化RabbitMQ连接池
func InitializeRabbitMQ() *Queue {
	return &Queue{
		connPool: sync.Pool{
			New: func() interface{} {
				conn, err := connectToRabbitMQ()
				if err != nil {
					App.Log.Error("Failed to create RabbitMQ connection: ", zap.Any("err", err))
				}
				return conn
			},
		},
	}
}

// connectToRabbitMQ 尝试连接到 RabbitMQ 集群
func connectToRabbitMQ() (*amqp091.Connection, error) {
	// rabbitURLs RabbitMQ 集群地址
	rabbitURLs := make([]string, 3)
	rabbitURLs = []string{
		"amqp://" + App.Config.Queue.Rabbitmq.Username + ":" + App.Config.Queue.Rabbitmq.Password + "@" + App.Config.Queue.Rabbitmq.Host + ":" + strconv.Itoa(App.Config.Queue.Rabbitmq.Port) + App.Config.Queue.Rabbitmq.Vhost,
		"amqp://" + App.Config.Queue.Rabbitmq.Username + ":" + App.Config.Queue.Rabbitmq.Password + "@" + App.Config.Queue.Rabbitmq.Host + ":" + strconv.Itoa(App.Config.Queue.Rabbitmq.Port) + App.Config.Queue.Rabbitmq.Vhost,
		"amqp://" + App.Config.Queue.Rabbitmq.Username + ":" + App.Config.Queue.Rabbitmq.Password + "@" + App.Config.Queue.Rabbitmq.Host + ":" + strconv.Itoa(App.Config.Queue.Rabbitmq.Port) + App.Config.Queue.Rabbitmq.Vhost,
	}

	var conn *amqp091.Connection
	var err error
	for _, url := range rabbitURLs {
		conn, err = amqp091.DialConfig(url, amqp091.Config{
			Heartbeat: 10 * time.Second,
			Locale:    "zh_CN.UTF-8",
		})
		if err == nil {
			App.Log.Info("Connected to RabbitMQ: " + url)
			return conn, nil
		}
		App.Log.Info("Failed to connect to RabbitMQ:" + url + ", retrying...")
		time.Sleep(2 * time.Second) // 重试间隔
	}
	return nil, err
}

// GetConnection 从连接池获取连接
func (qp *Queue) GetConnection() *amqp091.Connection {
	return qp.connPool.Get().(*amqp091.Connection)
}

// ReleaseConnection 将连接放回连接池
func (qp *Queue) ReleaseConnection(conn *amqp091.Connection) {
	qp.connPool.Put(conn)
}
