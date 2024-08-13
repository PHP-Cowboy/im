package initialize

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"im/global"
)

func InitRedis() {
	redisInfo := global.ServerConfig.RedisInfo

	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisInfo.Host, redisInfo.Port),
		Password: redisInfo.Password, // no password set
		DB:       redisInfo.Db,       // use default DB
		Username: redisInfo.Username,
	})

	global.Redis = &global.RedisCli{Cli: cli}
}
