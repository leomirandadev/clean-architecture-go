package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type memcacheImpl struct {
	cache *memcache.Client
	exp   time.Duration
}

func NewMemcache(opts Options) Cache {
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

func (c *memcacheImpl) Get(ctx context.Context, key string, v any) error {
	item, err := c.cache.Get(key)
	if err != nil {
		return err
	}

	return decode(item.Value, v)
}

func (c *memcacheImpl) Set(ctx context.Context, key string, v any) error {
	value, err := encode(v)
	if err != nil {
		return err
	}

	return c.cache.Set(&memcache.Item{Key: key, Value: value, Expiration: int32(c.exp)})
}

func (c *memcacheImpl) Del(ctx context.Context, key string) error {
	return c.cache.Delete(key)
}

func (c memcacheImpl) WithExpiration(d time.Duration) Cache {
	c.exp = d / time.Second
	return &c
}

func encode(v any) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func decode(data []byte, v any) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	return decoder.Decode(v)
}
