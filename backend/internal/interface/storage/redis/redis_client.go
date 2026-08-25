package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb       *redis.Client
	memStore  sync.Map
	isInMemory bool
}

type memoryEntry struct {
	val       string
	expiresAt time.Time
}

func NewRedisClient(addr, password string) *RedisClient {
	if addr != "" {
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       0,
		})
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := rdb.Ping(ctx).Err(); err == nil {
			return &RedisClient{rdb: rdb, isInMemory: false}
		}
	}

	return &RedisClient{isInMemory: true}
}

func (c *RedisClient) SetNX(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	if !c.isInMemory && c.rdb != nil {
		return c.rdb.SetNX(ctx, key, value, expiration).Result()
	}

	now := time.Now()
	val, loaded := c.memStore.LoadOrStore(key, memoryEntry{
		val:       value,
		expiresAt: now.Add(expiration),
	})

	if loaded {
		entry := val.(memoryEntry)
		if now.After(entry.expiresAt) {
			c.memStore.Store(key, memoryEntry{
				val:       value,
				expiresAt: now.Add(expiration),
			})
			return true, nil
		}
		return false, nil
	}

	return true, nil
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	if !c.isInMemory && c.rdb != nil {
		return c.rdb.Get(ctx, key).Result()
	}

	val, ok := c.memStore.Load(key)
	if !ok {
		return "", fmt.Errorf("redis: nil")
	}

	entry := val.(memoryEntry)
	if time.Now().After(entry.expiresAt) {
		c.memStore.Delete(key)
		return "", fmt.Errorf("redis: nil")
	}

	return entry.val, nil
}

func (c *RedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if !c.isInMemory && c.rdb != nil {
		return c.rdb.Set(ctx, key, value, expiration).Err()
	}

	c.memStore.Store(key, memoryEntry{
		val:       value,
		expiresAt: time.Now().Add(expiration),
	})
	return nil
}

func (c *RedisClient) Del(ctx context.Context, key string) error {
	if !c.isInMemory && c.rdb != nil {
		return c.rdb.Del(ctx, key).Err()
	}

	c.memStore.Delete(key)
	return nil
}
