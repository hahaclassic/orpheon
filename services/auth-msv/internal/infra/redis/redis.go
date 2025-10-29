package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/hahaclassic/orpheon/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisConfig = config.RedisConfig

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
