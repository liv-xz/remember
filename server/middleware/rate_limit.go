package middleware

import (
	"context"
	"net/http"
	"sync"

	"server/config"
	"server/database/redis"
	"server/errcode"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	authLimiterFallback    = make(map[string]*rate.Limiter)
	generalLimiterFallback = make(map[string]*rate.Limiter)
	fallbackMu             sync.RWMutex
)

func limiterFromPool(pool map[string]*rate.Limiter, ip string, r rate.Limit, burst int) *rate.Limiter {
	fallbackMu.Lock()
	defer fallbackMu.Unlock()

	if lim, ok := pool[ip]; ok {
		return lim
	}
	lim := rate.NewLimiter(r, burst)
	pool[ip] = lim
	return lim
}

var rateLimitScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCR", key)
if current == 1 then
    redis.call("EXPIRE", key, window)
end
if current > limit then
    return 0
end
return 1
`

func redisRateLimit(ip, prefix string, limitPerWindow int, windowSec int) bool {
	ctx := context.Background()
	key := prefix + ":" + ip
	result, err := redis.Client.Eval(ctx, rateLimitScript, []string{key}, limitPerWindow, windowSec).Int()
	if err != nil {
		return false
	}
	return result == 1
}

func rateLimitMiddleware(prefix string, r rate.Limit, burst int) gin.HandlerFunc {
	limitPerWindow := int(float64(r) * 60)
	if limitPerWindow <= 0 {
		limitPerWindow = 1
	}

	var fallbackPool map[string]*rate.Limiter
	switch prefix {
	case "ratelimit:auth":
		fallbackPool = authLimiterFallback
	default:
		fallbackPool = generalLimiterFallback
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if redis.IsAvailable() {
			if redisRateLimit(ip, prefix, limitPerWindow, 61) {
				c.Next()
				return
			}
		} else {
			limiter := limiterFromPool(fallbackPool, ip, r, burst)
			if limiter.Allow() {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
			"code":    errcode.RateLimited,
			"message": errcode.GetMessage(errcode.RateLimited),
		})
	}
}

func GeneralRateLimit() gin.HandlerFunc {
	cfg := config.Global.RateLimit
	if !cfg.Enable {
		return func(c *gin.Context) { c.Next() }
	}
	windowSec := 60
	limitPerMin := cfg.GeneralRate
	return rateLimitMiddleware("ratelimit:general", rate.Limit(limitPerMin)/rate.Limit(windowSec), cfg.GeneralBurst)
}

func AuthRateLimit() gin.HandlerFunc {
	cfg := config.Global.RateLimit
	if !cfg.Enable {
		return func(c *gin.Context) { c.Next() }
	}
	return rateLimitMiddleware("ratelimit:auth", rate.Limit(cfg.AuthRate)/60, cfg.AuthBurst)
}
