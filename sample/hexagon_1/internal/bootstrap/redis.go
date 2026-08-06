package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	sAddr := fmt.Sprintf("%s:%s", CONFIG.REDIS.HOST, CONFIG.REDIS.PORT)

	oClient := redis.NewClient(&redis.Options{
		Addr:     sAddr,
		Username: CONFIG.REDIS.USERNAME,
		Password: CONFIG.REDIS.PASSWORD,
		DB:       CONFIG.REDIS.DB,
	})

	if err := oClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return oClient, nil
}
