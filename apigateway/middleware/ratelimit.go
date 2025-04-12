package middleware

import (
	"net/http"
	"shortLink/common/ratelimit"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// 速率限流器实例
	rateLimiter ratelimit.RateLimiter
	// 批量限流器实例
	batchRateLimiter ratelimit.RateLimiter
	// 初始化互斥锁
	once sync.Once
)

// 初始化限流器
func initLimiters() {
	once.Do(func() {
		// 初始化全局限流器
		globalCfg := &ratelimit.Config{
			Type:        "token_bucket",
			WindowSize:  time.Second,
			MaxRequests: 100,
			Rate:        100,
			Capacity:    100,
		}
		var err error
		rateLimiter, err = ratelimit.NewRateLimiter(globalCfg)
		if err != nil {
			panic(err)
		}

		// 初始化批量操作限流器
		batchCfg := &ratelimit.Config{
			Type:        "sliding_window",
			WindowSize:  time.Minute,
			MaxRequests: 10,
		}
		batchRateLimiter, err = ratelimit.NewRateLimiter(batchCfg)
		if err != nil {
			panic(err)
		}
	})
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware() gin.HandlerFunc {
	// 确保限流器已初始化
	initLimiters()

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !rateLimiter.Allow(c.Request.Context(), clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// BatchRateLimitMiddleware 批量创建短链接的限流中间件
func BatchRateLimitMiddleware() gin.HandlerFunc {
	// 确保限流器已初始化
	initLimiters()

	return func(c *gin.Context) {
		userID := strconv.Itoa(int(c.GetUint("UserID")))

		if !batchRateLimiter.Allow(c.Request.Context(), userID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "批量创建请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
