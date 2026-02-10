package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheStore interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, interface{}, time.Duration) error
}
type RedisCacheStore struct {
	client *redis.Client
}

func NewRedisCacheStore(client *redis.Client) *RedisCacheStore {
	return &RedisCacheStore{
		client: client,
	}

}

func (c *RedisCacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCacheStore) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}
