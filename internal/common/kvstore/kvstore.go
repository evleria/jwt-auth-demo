package kvstore

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type KVStore interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) (int64, error)
}

type redisStore struct {
	ctx    context.Context
	client *redis.Client
}

func New(ctx context.Context, opts *redis.Options) (KVStore, error) {
	client := redis.NewClient(opts)
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &redisStore{
		ctx:    ctx,
		client: client,
	}, nil
}

func (c *redisStore) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

func (c *redisStore) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

func (c *redisStore) Delete(key string) (int64, error) {
	return c.client.Del(c.ctx, key).Result()
}
