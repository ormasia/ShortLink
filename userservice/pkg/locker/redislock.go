package locker

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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
	return &RedisLock{
		client: client,
		key:    key,
		value:  uuid.NewString(), // 随机唯一值，确保是当前持有锁者
		expire: expire,
		ctx:    ctx,
		cancel: cancel,
	}
}

// 尝试加锁（非阻塞）
func (l *RedisLock) TryLock() (bool, error) {
	ok, err := l.client.SetNX(l.ctx, l.key, l.value, l.expire).Result()
	if err != nil {
		return false, err
	}
	l.acquired = ok
	return ok, nil
}

// 解锁（只释放自己加的锁）
func (l *RedisLock) Unlock() error {
	defer l.cancel()

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
		return err
	}
	if res.(int64) == 0 {
		return errors.New("解锁失败: 不是持有者")
	}
	return nil
}
