package main

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Connect(t *testing.T) {
	rc := &RedisCache{}

	// Mock Redis client
	mockClient, mock := redismock.NewClientMock()
	mock.ExpectPing().SetVal("PONG")

	rc.rdb = mockClient
	err := rc.Connect("localhost", "6379")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisCache_GetCachedMessages(t *testing.T) {
	rc := &RedisCache{}

	// Mock Redis client
	mockClient, mock := redismock.NewClientMock()
	mock.ExpectGet("chat:test:page:1").SetVal("cached_message")

	rc.rdb = mockClient
	result, err := rc.GetCachedMessages("test")

	assert.NoError(t, err)
	assert.Equal(t, "cached_message", result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisCache_CacheMessages(t *testing.T) {
	rc := &RedisCache{}

	messages := []Message{{Content: "cached_message"}}
	// Mock Redis client
	mockClient, mock := redismock.NewClientMock()
	mock.ExpectSet("chat:test:page:1", messages, 24*time.Hour).SetVal("OK")

	rc.rdb = mockClient
	err := rc.CacheMessages("test", messages)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
