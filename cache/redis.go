package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func InitRedis(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func Set(key, value string) {
	rdb.Set(ctx, key, value, time.Hour*24).Err()
}

func Get(key string) string {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}
