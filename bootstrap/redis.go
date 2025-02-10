package bootstrap

import (
	"context"
	"shalabing-gin/global"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func InitializeRedis() *redis.ClusterClient { // *redis.Client
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     global.App.Config.Redis.Host + ":" + strconv.Itoa(global.App.Config.Redis.Port),
	// 	Password: global.App.Config.Redis.Password, // no password set
	// 	DB:       global.App.Config.Redis.DB,       // use default DB
	// })
	// _, err := client.Ping(context.Background()).Result()
	// if err != nil {
	// 	global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
	// 	return nil
	// }
	// return client

	/*
		数据库编号:
		在 Redis 集群中，每个节点默认只有一个数据库（DB 0），不支持使用多个数据库。
		因此，数据库编号（如 DB 参数）在集群模式下被忽略。
		连接复用 go-redis 内部已经实现了高效的连接复用，设置合适的参数可以在高并发场景下发挥更好的性能。
	*/

	redis := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        global.App.Config.Redis.Addrs,
		Password:     global.App.Config.Redis.Password,     // 设置 Redis 密码
		PoolSize:     global.App.Config.Redis.PoolSize,     // 限制最大连接数，防止 Redis 集群压力过大。
		MinIdleConns: 5,                                    // 设置最小空闲连接数，保证在高并发情况下随时有连接可用。
		PoolTimeout:  5 * time.Second,                      // 当连接池满时，等待可用连接的超时时间。
		DialTimeout:  global.App.Config.Redis.DialTimeout,  // 连接建立的最大超时时间。
		ReadTimeout:  global.App.Config.Redis.ReadTimeout,  // 读取超时时间,避免长时间阻塞。
		WriteTimeout: global.App.Config.Redis.WriteTimeout, // 写入超时时间,避免长时间阻塞。
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redis.Ping(ctx).Err(); err != nil {
		global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
	}
	return redis
}
