package services_common

import (
	"context"
	"fmt"
	"shalabing-gin/global"
	"time"
)

// 全局使用,就需要把定义成公有的
var ctxRedis = context.Background()

const CAPTCHA = "captcha:"

type RedisStore struct {
}

// 实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	//time.Minute*2：有效时间2分钟
	err := global.App.Redis.Set(ctxRedis, key, value, time.Minute*60).Err()

	return err
}

// 实现获取captcha的方法
func (r RedisStore) Get(id string, clear bool) string {
	key := CAPTCHA + id
	val, err := global.App.Redis.Get(ctxRedis, key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		err := global.App.Redis.Del(ctxRedis, key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

// 实现验证captcha的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {
	v := RedisStore{}.Get(id, clear)
	//fmt.Println("key:"+id+";value:"+v+";answer:"+answer)
	return v == answer
}
