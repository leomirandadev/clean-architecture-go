package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/leomirandadev/clean-architecture-go/utils/logger"
)

type memcacheImpl struct {
	cache *memcache.Client
	exp   time.Duration
}

// NewMemcache cria uma implementação da interface Cache para utilizar o memcache
func NewMemcache(opts Options, log logger.Logger) Cache {
	client := memcache.New(opts.URL)
	err := client.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &memcacheImpl{
		exp:   opts.Expiration,
		cache: client,
	}
}

func (c *memcacheImpl) Get(ctx context.Context, key string, v interface{}) bool {
	item, err := c.cache.Get(key)
	if err != nil {
		return false
	}

	if err := decode(item.Value, v); err != nil {
		return false
	}

	return true
}

func (c *memcacheImpl) Set(ctx context.Context, key string, v interface{}) bool {
	value, err := encode(v)
	if err != nil {
		return false
	}

	err = c.cache.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(c.exp)})
	if err != nil {
		return false
	}

	return true
}

func (c *memcacheImpl) Del(ctx context.Context, key string) bool {
	err := c.cache.Delete(key)
	if err != nil {
		return false
	}

	return true
}

func (c memcacheImpl) WithExpiration(d time.Duration) Cache {
	c.exp = d / time.Second
	return &c
}

func encode(v interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decode(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	return decoder.Decode(v)
}
