package circuitbreaker

import (
	"context"
	"errors"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker 定义断路器接口
type CircuitBreaker interface {
	Execute(ctx context.Context, req interface{}, operation func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error)
}

// DefaultCircuitBreaker 默认断路器实现
type DefaultCircuitBreaker struct {
	cb *gobreaker.CircuitBreaker
}

// NewDefaultCircuitBreaker 创建默认断路器
func NewDefaultCircuitBreaker(name string) *DefaultCircuitBreaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,                // 半开状态下的最大请求数
		Interval:    10 * time.Second, // 重置计数器的间隔
		Timeout:     5 * time.Second,  // 断路器打开状态的超时时间
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6 // 请求数>=3且失败率>=60%时触发断路器
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// 状态变化时的回调
		},
	}

	return &DefaultCircuitBreaker{
		cb: gobreaker.NewCircuitBreaker(settings),
	}
}

// Execute 执行操作，如果断路器打开则返回错误
func (cb *DefaultCircuitBreaker) Execute(ctx context.Context, req interface{}, operation func(ctx context.Context, req interface{}) (interface{}, error)) (interface{}, error) {
	// 检查上下文是否已取消
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// 使用断路器执行操作
	result, err := cb.cb.Execute(func() (interface{}, error) {
		return operation(ctx, req)
	})

	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) {
			// 断路器打开状态
			return nil, errors.New("服务暂时不可用，请稍后重试")
		}
		if errors.Is(err, gobreaker.ErrTooManyRequests) {
			// 半开状态下请求过多
			return nil, errors.New("服务正在恢复中，请稍后重试")
		}
	}

	return result, err
}
