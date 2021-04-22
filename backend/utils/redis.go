package utils

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/batijo/video-platform/backend/models"
	redis "github.com/go-redis/redis/v8"
)

var (
	RedisCl *redis.Client
	ctx     = context.Background()
)

func InitRedisClient(port string) {
	RedisCl = redis.NewClient(&redis.Options{
		Addr:     "redis:" + port,
		Password: "super_secret_123",
		DB:       0, // use default DB
	})

	_, err := RedisCl.Ping(ctx).Result()
	if err != nil {
		log.Panicln(err)
	}
}

func CreateRedisAuth(userID uint, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	now := time.Now()

	err := RedisCl.Set(ctx, td.AccessUuid, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func DeleteRedisAuth(tokenUuid string) (int64, error) {
	deleted, err := RedisCl.Del(ctx, tokenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func fetchAuth(ad *models.AccessDetails) (uint, error) {
	uid, err := RedisCl.Get(ctx, ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, err := strconv.ParseUint(uid, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}
