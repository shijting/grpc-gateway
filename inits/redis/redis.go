package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/showiot/camera/inits/config"
)
var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init() (err error) {
	redisConf := config.Conf.RedisConfig
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password:     redisConf.Password,
		DB:           redisConf.DB,
		PoolSize:     redisConf.PoolSize,
		MinIdleConns: 3,
	})

	_, err = client.Ping().Result()
	return
}

func GetDB() *redis.Client {
	return client
}
func Close() (err error){
	if client !=nil {
		err = client.Close()
	}
	return
}