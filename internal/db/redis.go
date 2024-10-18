package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-log-keeper/config"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0, //default redisDB,
	})
}
