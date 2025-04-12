package ratelimit

import (
	"context"
	"math"
	"sync"
	"time"
)

// tokenBucketLimiter 令牌桶限流器
type tokenBucketLimiter struct {
	rate       float64    // 令牌产生速率
	capacity   int64      // 桶的容量
	tokens     float64    // 当前令牌数
	lastRefill time.Time  // 上次填充时间
	mu         sync.Mutex // 互斥锁
}

// NewTokenBucketLimiter 创建新的令牌桶限流器
func NewTokenBucketLimiter(cfg *Config) (RateLimiter, error) {
	return &tokenBucketLimiter{
		rate:       cfg.Rate,
		capacity:   cfg.Capacity,
		tokens:     float64(cfg.Capacity),
		lastRefill: time.Now(),
	}, nil
}

// Allow 检查是否允许请求通过
func (l *tokenBucketLimiter) Allow(ctx context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 计算从上次填充到现在产生的令牌数
	now := time.Now()
	elapsed := now.Sub(l.lastRefill).Seconds()
	newTokens := elapsed * l.rate

	// 更新令牌数和上次填充时间
	l.tokens = math.Min(float64(l.capacity), l.tokens+newTokens)
	l.lastRefill = now

	// 如果有足够的令牌，允许请求通过
	if l.tokens >= 1 {
		l.tokens--
		return true
	}

	return false
}

// Close 关闭限流器
func (l *tokenBucketLimiter) Close() error {
	return nil
}
