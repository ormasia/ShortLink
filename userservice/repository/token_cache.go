package repository

import (
	"context"
	"fmt"
	"shortLink/userservice/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// TokenCache 定义缓存接口
type TokenCache interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

// RedisCache 实现基于Redis的缓存
type RedisCache struct {
	rdb        *redis.Client
	ctx        context.Context
	expiration time.Duration
}

// NewRedisCache 创建一个新的RedisCache实例
func NewRedisCache(rdb *redis.Client, expiration time.Duration) *RedisCache {
	return &RedisCache{
		rdb:        rdb,
		ctx:        context.Background(),
		expiration: expiration,
	}
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value string) error {
	if c.rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}
	err := c.rdb.Set(c.ctx, key, value, c.expiration).Err()
	if err != nil {
		logger.Log.Error("设置缓存失败", zap.Error(err))
		return err
	}
	logger.Log.Debug("设置缓存成功", zap.String("key", key))
	return nil
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (string, error) {
	if c.rdb == nil {
		return "", fmt.Errorf("Redis未初始化")
	}
	val, err := c.rdb.Get(c.ctx, key).Result()
	if err != nil {
		logger.Log.Debug("获取缓存失败", zap.Error(err), zap.String("key", key))
		return "", err
	}
	logger.Log.Debug("获取缓存成功", zap.String("key", key))
	return val, nil
}

// Delete 删除缓存
func (c *RedisCache) Delete(key string) error {
	if c.rdb == nil {
		return fmt.Errorf("Redis未初始化")
	}
	err := c.rdb.Del(c.ctx, key).Err()
	if err != nil {
		logger.Log.Error("删除缓存失败", zap.Error(err), zap.String("key", key))
		return err
	}
	logger.Log.Debug("删除缓存成功", zap.String("key", key))
	return nil
}
