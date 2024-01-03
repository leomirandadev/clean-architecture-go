//go:generate mockgen -source cache.go -destination mocks/cache_mock.go -package mocks
package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, v any) error
	Set(ctx context.Context, key string, v any) error
	Del(ctx context.Context, key string) error
	WithExpiration(d time.Duration) Cache
}

type Options struct {
	Expiration time.Duration `mapstructure:"expiration"`
	URL        string        `mapstructure:"url"`
}
