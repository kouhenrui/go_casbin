package test

import (
	"context"
	"go_casbin/pkg/redis"
	"testing"
)

func TestRedis(t *testing.T) {
	err := redis.InitRedis(redis.RedisOptions{
		Addr: "localhost:6379",
	})
	if err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}
	client:=redis.GetRedisInstance()
	client.Set(context.Background(), "test", "test", 0)
}