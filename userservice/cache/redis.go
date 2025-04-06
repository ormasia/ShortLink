package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
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
	// 测试连接
	if err := rdb.Ping(ctx).Err(); err != nil {
		fmt.Printf("Redis连接失败: %v", err)
	}
	fmt.Println("Redis连接成功")
}

// 获取redis
func GetRedis() *redis.Client {
	return rdb
}

func Set(key, value string) {
	if rdb == nil {
		fmt.Println("Redis未初始化")
		return
	}
	err := rdb.Set(ctx, key, value, time.Hour*24).Err()
	if err != nil {
		fmt.Printf("设置缓存失败: %v", err)
	}
	fmt.Println("设置缓存成功")
}

func Get(key string) string {
	if rdb == nil {
		fmt.Println("Redis未初始化")
		return ""
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	fmt.Println("获取缓存成功", val)
	return val
}
