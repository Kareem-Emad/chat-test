package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	GetCachedMessages(chat string) (string, error)
	CacheMessages(chat string, messages []Message) error
	Connect(host, port string) error
}

type RedisCache struct {
	rdb *redis.Client
}

func (rc *RedisCache) Connect(host, port string) error {
	// Connect to Redis
	if rc.rdb == nil {
		rc.rdb = redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		})
	}
	fmt.Println("cache connected")

	return rc.rdb.Ping(context.Background()).Err()
}

func (rc *RedisCache) GetCachedMessages(chat string) (string, error) {
	// Get cached messages from Redis
	cacheKey := fmt.Sprintf("chat:%s:page:1", chat)
	cachedMessages, err := rc.rdb.Get(context.Background(), cacheKey).Result()
	if err != nil {
		return "", err
	}

	return cachedMessages, nil
}

func (rc *RedisCache) CacheMessages(chat string, messages []Message) error {
	// Cache messages in Redis
	cacheKey := fmt.Sprintf("chat:%s:page:1", chat)
	res := rc.rdb.Set(context.Background(), cacheKey, messages, 24*time.Hour)
	return res.Err()
}
