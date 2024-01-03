package cache

import (
	"context"
	"encoding/json"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type redisImpl struct {
	cache *redis.Client
	exp   time.Duration
}

func NewRedis(opts Options) Cache {
	return &redisImpl{
		exp: opts.Expiration,
		cache: redis.NewClient(&redis.Options{
			Addr: opts.URL,
		}),
	}
}

func (c *redisImpl) Get(ctx context.Context, key string, v any) error {
	cmd := c.cache.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		return err
	}

	return cmd.Scan(v)
}

func (c *redisImpl) Set(ctx context.Context, key string, v any) error {
	valueByte, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.cache.Set(ctx, key, valueByte, c.exp).Err()
}

func (c *redisImpl) Del(ctx context.Context, key string) error {
	return c.cache.Del(ctx, key).Err()
}

func (c redisImpl) WithExpiration(d time.Duration) Cache {
	c.exp = d
	return &c
}
