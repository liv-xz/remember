package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"server/database/redis"
)

type LoginProtector struct {
	maxFailures  int
	lockDuration time.Duration

	mu      sync.RWMutex
	records map[string]*failRecord
}

type failRecord struct {
	failures  int
	lockUntil time.Time
}

func NewLoginProtector(maxFailures int, lockDuration time.Duration) *LoginProtector {
	return &LoginProtector{
		maxFailures:  maxFailures,
		lockDuration: lockDuration,
		records:      make(map[string]*failRecord),
	}
}

func (lp *LoginProtector) key(ip, username string) string {
	return fmt.Sprintf("login:fail:%s:%s", ip, username)
}

func (lp *LoginProtector) IsLocked(ip, username string) bool {
	k := lp.key(ip, username)

	if redis.IsAvailable() {
		ctx := context.Background()
		val, err := redis.Client.Get(ctx, k).Int()
		if err != nil {
			return false
		}
		return val >= lp.maxFailures
	}

	lp.mu.RLock()
	defer lp.mu.RUnlock()
	rec, ok := lp.records[k]
	if !ok || rec.lockUntil.IsZero() {
		return false
	}
	return time.Now().Before(rec.lockUntil)
}

func (lp *LoginProtector) LockRemaining(ip, username string) time.Duration {
	k := lp.key(ip, username)

	if redis.IsAvailable() {
		ctx := context.Background()
		ttl, err := redis.Client.TTL(ctx, k).Result()
		if err != nil || ttl <= 0 {
			return 0
		}
		return ttl.Truncate(time.Second)
	}

	lp.mu.RLock()
	defer lp.mu.RUnlock()
	rec, ok := lp.records[k]
	if !ok || rec.lockUntil.IsZero() {
		return 0
	}
	remaining := time.Until(rec.lockUntil)
	if remaining < 0 {
		return 0
	}
	return remaining.Truncate(time.Second)
}

func (lp *LoginProtector) RecordFailure(ip, username string) {
	k := lp.key(ip, username)

	if redis.IsAvailable() {
		ctx := context.Background()
		_, _ = redis.Client.Incr(ctx, k).Result()
		val, _ := redis.Client.Get(ctx, k).Int64()
		if val == 1 {
			bufferSec := int(lp.lockDuration.Seconds()) * 2
			redis.Client.Expire(ctx, k, time.Duration(bufferSec)*time.Second)
		}
		return
	}

	lp.mu.Lock()
	defer lp.mu.Unlock()
	rec, ok := lp.records[k]
	if !ok {
		rec = &failRecord{}
		lp.records[k] = rec
	}
	rec.failures++
	if rec.failures >= lp.maxFailures {
		rec.lockUntil = time.Now().Add(lp.lockDuration)
	}
}

func (lp *LoginProtector) Reset(ip, username string) {
	k := lp.key(ip, username)

	if redis.IsAvailable() {
		ctx := context.Background()
		redis.Client.Del(ctx, k)
		return
	}

	lp.mu.Lock()
	defer lp.mu.Unlock()
	delete(lp.records, k)
}
