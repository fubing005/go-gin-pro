package bootstrap

import (
	"shalabing-gin/global"
)

func InitializeService() {
	// 初始化配置
	InitializeConfig()

	// 初始化日志
	global.App.Log = InitializeLog()

	// 初始化数据库
	global.App.DB = InitializeDB()

	// 初始化Redis
	global.App.Redis = InitializeRedis()

	// 初始化MongoDB
	global.App.MongoDB = InitializeMongoDB()

	//初始化clickhouse
	global.App.Clickhouse = InitializeClickHouse()

	// 初始化ES
	global.App.ES = InitializeES()

	// 初始化队列
	global.App.Rabbitmq = global.InitializeRabbitMQ()
	global.App.Kafka = global.InitProducer()

	// 初始化验证器
	InitializeValidator()

	// 初始化文件系统
	InitializeStorage()

	// 初始化计划任务
	InitializeCron()
}
