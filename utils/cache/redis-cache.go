package cache

import (
	"context"
	"time"
	"encoding/json"

	redis "github.com/go-redis/redis/v8"
)

type redisImpl struct {
	cache *redis.Client
	exp   time.Duration
}

// NewRedis cria uma implementação da interface Cache para utilizar o redis
func NewRedis(opts Options) Cache {
	return &redisImpl{
		exp: opts.Expiration,
		cache: redis.NewClient(&redis.Options{
			Addr: opts.URL,
		}),
	}
}

func (c *redisImpl) Get(ctx context.Context, key string, v interface{}) bool {
	cmd := c.cache.Get(ctx, key)
	if cmd.Err() != nil {
		return false
	}

	if err := cmd.Scan(v); err != nil {
		return false
	}

	return true
}

func (c *redisImpl) Set(ctx context.Context, key string, v interface{}) bool {
	valueByte, err := json.Marshal(v)
	if err != nil {
		return false
	}
	cmd := c.cache.Set(ctx, key, valueByte, c.exp)
	if cmd.Err() != nil {
		return false
	}

	return true
}

func (c *redisImpl) Del(ctx context.Context, key string) bool {
	cmd := c.cache.Del(ctx, key)
	if cmd.Err() != nil {
		return false
	}

	return true
}

func (c redisImpl) WithExpiration(d time.Duration) Cache {
	c.exp = d
	return &c
}
