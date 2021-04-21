package utils

import (
	"context"
	"log"

	redis "github.com/go-redis/redis/v8"
)

var RedisCl *redis.Client

func InitRedisClient() {
	RedisCl := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "super_secret_123",
		DB:       0, // use default DB
	})

	_, err := RedisCl.Ping(context.Background()).Result()
	if err != nil {
		log.Panicln(err)
	}

}
