package global

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
)

type Interface interface {
	Get() bool
	Block(seconds int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string // 锁名称
	owner   string // 锁标识
	seconds int64  // 有效期
}

// 释放锁 Lua 脚本，防止任何客户端都能解锁
const releaseLockLuaScript = `
if redis.call("get",KEYS[1]) == ARGV[1] then
    return redis.call("del",KEYS[1])
else
    return 0
end
`

// 生成锁
func Lock(name string, seconds int64) Interface {
	return &lock{
		context.Background(),
		name,
		RandString(16),
		seconds,
	}
}

// 获取锁
func (l *lock) Get() bool {
	return App.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// 阻塞一段时间，尝试获取锁
func (l *lock) Block(seconds int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-seconds >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

// 释放锁
// func (l *lock) Release() bool {
// 	luaScript := redis.NewScript(releaseLockLuaScript)
// 	result := luaScript.Run(l.context, App.Redis, []string{l.name}, l.owner).Val().(int64)
// 	return result != 0
// }

// 释放锁
func (l *lock) Release() bool {
	// 创建Lua脚本对象
	luaScript := redis.NewScript(releaseLockLuaScript)
	// 确保键在同一个哈希槽中，使用哈希标签
	hashedName := "{" + l.name + "}"
	// 执行Lua脚本
	result, err := luaScript.Run(l.context, App.Redis, []string{hashedName}, l.owner).Int64()
	if err != nil {
		// 处理错误
		return false
	}
	return result != 0
}

// 强制释放锁
func (l *lock) ForceRelease() {
	App.Redis.Del(l.context, l.name).Val()
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
