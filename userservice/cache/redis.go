package cache

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

// 初始化redis
func InitRedis(addr, password string, port, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", addr, port),
		Password: password,
		DB:       db,
	})
	// // 测试连接
	// if err := rdb.Ping(ctx).Err(); err != nil {
	// 	logger.Log.Error("Redis连接失败", zap.Error(err))
	// }
	// logger.Log.Info("Redis连接成功")
}

// 获取redis
func GetRedis() *redis.Client {
	return rdb
}
