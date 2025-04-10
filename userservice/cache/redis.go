package cache

import (
	"context"
	"fmt"
	"time"

	"shortLink/userservice/logger"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	ctx = context.Background()
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

func Set(key, value string) {
	if rdb == nil {
		logger.Log.Warn("Redis未初始化")
		return
	}
	err := rdb.Set(ctx, key, value, time.Hour*24).Err()
	if err != nil {
		logger.Log.Error("设置缓存失败", zap.Error(err))
	}
	logger.Log.Info("设置缓存成功")
}

func Get(key string) string {
	if rdb == nil {
		logger.Log.Warn("Redis未初始化")
		return ""
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		logger.Log.Debug("获取缓存失败", zap.Error(err), zap.String("key", key))
		return ""
	}
	logger.Log.Info("获取缓存成功", zap.String("value", val))
	return val
}

func Del(key string) {
	if rdb == nil {
		logger.Log.Warn("Redis未初始化")
		return
	}
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		logger.Log.Error("删除缓存失败", zap.Error(err))
	}
}
