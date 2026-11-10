package redis

import (
	"context"
	"time"
)

func BlacklistToken(hash string, ttl time.Duration) {
	if Client == nil {
		return
	}
	Client.Set(context.Background(), "blacklist:"+hash, 1, ttl)
}

func IsTokenBlacklisted(hash string) bool {
	if Client == nil {
		return false
	}
	_, err := Client.Get(context.Background(), "blacklist:"+hash).Result()
	return err == nil
}
