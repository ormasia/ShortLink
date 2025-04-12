package ratelimit

import (
	"context"
	"time"
)

// RateLimiter 定义限流器接口
type RateLimiter interface {
	// Allow 检查是否允许请求通过
	Allow(ctx context.Context, key string) bool
	// Close 关闭限流器
	Close() error
}

// Config 限流器配置
type Config struct {
	// 限流类型：token_bucket, sliding_window
	Type string
	// 时间窗口大小
	WindowSize time.Duration
	// 每个时间窗口允许的请求数
	MaxRequests int64
	// 令牌桶速率（每秒）
	Rate float64
	// 令牌桶容量
	Capacity int64
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(cfg *Config) (RateLimiter, error) {
	switch cfg.Type {
	case "token_bucket":
		return NewTokenBucketLimiter(cfg)
	case "sliding_window":
		return NewSlidingWindowLimiter(cfg)
	default:
		return NewTokenBucketLimiter(cfg)
	}
}
