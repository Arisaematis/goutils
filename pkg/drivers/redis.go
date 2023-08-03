package drivers

import (
	"context"
	"github.com/go-redis/redis/v8"
	"goutils/pkg/setting"
)

// RedisDB 定义一个全局变量
var RedisDB *redis.Client

func InitRedis() {
	RedisDB = redis.NewClient(&redis.Options{
		Addr:     setting.RedisSetting.RedisHost, // 指定 redis 连接地址
		Password: setting.RedisSetting.RedisPwd,
		DB:       setting.RedisSetting.RedisDB, // redis一共16个库，指定其中一个库即可
	})
	_, err := RedisDB.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
