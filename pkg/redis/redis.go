package redis

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once   sync.Once
)

type RedisService interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, keys ...string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	HSet(ctx context.Context, key string, values ...interface{}) error
	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	LPush(ctx context.Context, key string, values ...interface{}) error
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Publish(ctx context.Context, channel string, message interface{}) error
	Subscribe(ctx context.Context, channel string) *redis.PubSub
	TryLock(ctx context.Context, key, value string, expiration time.Duration) (bool, error)
	Unlock(ctx context.Context, key, value string) (bool, error)
}
type RedisServiceImpl struct {
	client *redis.Client
}
type RedisOptions struct {
	Addr     string  // 地址
	Username *string // 用户名
	Password *string // 密码
}
func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
func InitRedis(options RedisOptions)error {
	var err error
	once.Do(func() {
		client = redis.NewClient(&redis.Options{
			Addr:     options.Addr,
			Username: getString(options.Username),
			Password: getString(options.Password),
		})
		_, err = client.Ping(context.Background()).Result()
	})
	return err
}
func GetRedisInstance() RedisServiceImpl {
	return RedisServiceImpl{client: client}
}

// Set
func (r *RedisServiceImpl) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get
func (r *RedisServiceImpl) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Del
func (r *RedisServiceImpl) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists
func (r *RedisServiceImpl) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(ctx, keys...).Result()
}

// Expire
func (r *RedisServiceImpl) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

// Incr
func (r *RedisServiceImpl) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Decr
func (r *RedisServiceImpl) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(ctx, key).Result()
}

// HSet
func (r *RedisServiceImpl) HSet(ctx context.Context, key string, values ...interface{}) error {
	return r.client.HSet(ctx, key, values...).Err()
}

// HGet
func (r *RedisServiceImpl) HGet(ctx context.Context, key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

// HGetAll
func (r *RedisServiceImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

// LPush
func (r *RedisServiceImpl) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.client.LPush(ctx, key, values...).Err()
}

// RPop
func (r *RedisServiceImpl) RPop(ctx context.Context, key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

// LRange
func (r *RedisServiceImpl) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

// Publish
func (r *RedisServiceImpl) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.client.Publish(ctx, channel, message).Err()
}

// Subscribe
func (r *RedisServiceImpl) Subscribe(ctx context.Context, channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

// TryLock 尝试加锁
func (r *RedisServiceImpl) TryLock(ctx context.Context, key, value string, expiration time.Duration) (bool, error) {
	ok, err := r.client.SetNX(ctx, key, value, expiration).Result()
	return ok, err
}

// Unlock 释放锁（只删除自己加的锁，防止误删）
func (r *RedisServiceImpl) Unlock(ctx context.Context, key, value string) (bool, error) {
	// 用 Lua 脚本保证原子性
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`
	res, err := r.client.Eval(ctx, script, []string{key}, value).Result()
	return res.(int64) == 1, err
}