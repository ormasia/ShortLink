package ratelimit

import (
	"context"
	"sync"
	"time"
)

// slidingWindowLimiter 滑动窗口限流器
type slidingWindowLimiter struct {
	windowSize  time.Duration // 窗口大小
	maxRequests int64         // 窗口内最大请求数
	requests    []time.Time   // 请求时间列表
	mu          sync.Mutex    // 互斥锁
}

// NewSlidingWindowLimiter 创建新的滑动窗口限流器
func NewSlidingWindowLimiter(cfg *Config) (RateLimiter, error) {
	return &slidingWindowLimiter{
		windowSize:  cfg.WindowSize,
		maxRequests: cfg.MaxRequests,
		requests:    make([]time.Time, 0, cfg.MaxRequests),
	}, nil
}

// Allow 检查是否允许请求通过
func (l *slidingWindowLimiter) Allow(ctx context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.windowSize)

	// 移除窗口外的请求
	for len(l.requests) > 0 && l.requests[0].Before(windowStart) {
		l.requests = l.requests[1:]
	}

	// 检查是否超过最大请求数
	if int64(len(l.requests)) >= l.maxRequests {
		return false
	}

	// 添加新请求
	l.requests = append(l.requests, now)
	return true
}

// Close 关闭限流器
func (l *slidingWindowLimiter) Close() error {
	return nil
}
