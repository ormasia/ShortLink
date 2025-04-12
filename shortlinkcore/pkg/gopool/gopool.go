package gopool

import (
	"sync"

	"github.com/panjf2000/ants/v2"
)

var pool *ants.Pool
var once sync.Once

func GetPool() *ants.Pool {
	return pool
}

func init() {
	once.Do(func() {
		pool, _ = ants.NewPool(20) // 最多20个并发检测任务
	})
}
