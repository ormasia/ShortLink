package locker

import (
	"context"
	"errors"
	"time"

	"shortLink/shortlinkcore/logger"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type RedisLock struct {
	client   *redis.Client
	key      string
	value    string
	expire   time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
	acquired bool
}

// 创建一个 Redis 分布式锁对象
func NewRedisLock(client *redis.Client, key string, expire time.Duration) *RedisLock {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	value := uuid.NewString()
	logger.Log.Debug("创建Redis分布式锁",
		zap.String("key", key),
		zap.String("value", value),
		zap.Duration("expire", expire))
	return &RedisLock{
		client: client,
		key:    key,
		value:  value, // 随机唯一值，确保是当前持有锁者
		expire: expire,
		ctx:    ctx,
		cancel: cancel,
	}
}

// 尝试加锁（非阻塞）
func (l *RedisLock) TryLock() (bool, error) {
	ok, err := l.client.SetNX(l.ctx, l.key, l.value, l.expire).Result()
	if err != nil {
		logger.Log.Error("Redis锁获取失败",
			zap.String("key", l.key),
			zap.Error(err))
		return false, err
	}
	l.acquired = ok
	if ok {
		logger.Log.Debug("Redis锁获取成功",
			zap.String("key", l.key),
			zap.String("value", l.value),
			zap.Duration("expire", l.expire))
	} else {
		logger.Log.Debug("Redis锁已被占用",
			zap.String("key", l.key))
	}
	return ok, nil
}

// 解锁（只释放自己加的锁）
func (l *RedisLock) Unlock() error {
	defer l.cancel()

	logger.Log.Debug("尝试释放Redis锁",
		zap.String("key", l.key),
		zap.String("value", l.value))

	// Lua 脚本保证原子性：只有 value 一致才删除
	const script = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
else
	return 0
end
`
	res, err := l.client.Eval(l.ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		logger.Log.Error("Redis锁释放失败",
			zap.String("key", l.key),
			zap.Error(err))
		return err
	}
	if res.(int64) == 0 {
		logger.Log.Warn("Redis锁释放失败: 不是持有者",
			zap.String("key", l.key),
			zap.String("value", l.value))
		return errors.New("解锁失败: 不是持有者")
	}
	logger.Log.Debug("Redis锁释放成功",
		zap.String("key", l.key))
	return nil
}
