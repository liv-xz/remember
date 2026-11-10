package redis

import (
	"context"
	"fmt"
	"time"

	"server/configmodel"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init(cfg configmodel.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		Client = nil
		return fmt.Errorf("Redis 连接失败: %w", err)
	}
	return nil
}

func IsAvailable() bool {
	if Client == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return Client.Ping(ctx).Err() == nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
