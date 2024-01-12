package cache

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	ErrNotFound         = CacheError("[cache] not found")
	ErrUnsuportedSchema = CacheError("[cache] unsupported scheme")
)

type CacheError string

func (e CacheError) Error() string {
	return string(e)
}

type (
	InitFunc      func(netURL *url.URL) (Cache, error)
	DeleteOptions func(options *DeleteCache)
)

type DeleteCache struct {
	Pattern string
}

type CacheImpl struct {
	client *redis.Client
	ns     string
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Get(ctx context.Context, key string) ([]byte, error)
	GetObject(ctx context.Context, key string, doc interface{}) error
	GetString(ctx context.Context, key string) (string, error)
	GetInt(ctx context.Context, key string) (int64, error)
	GetFloat(ctx context.Context, key string) (float64, error)
	Exist(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string, opts ...DeleteOptions) error
	GetKeys(ctx context.Context, pattern string) []string
	RemainingTime(ctx context.Context, key string) int
	Publish(ctx context.Context, channel, message string) error
	Subscribe(ctx context.Context, topic string) (Subscriber, error)
	Close() error
}

func NewCacheRepository(client *redis.Client, ns string) Cache {
	return &CacheImpl{
		client: client,
		ns:     ns}
}

func (c *CacheImpl) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	switch value.(type) {
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, []byte:
		return c.client.
			Set(ctx, c.ns+key, value, time.Duration(expiration)*time.Second).
			Err()
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return c.client.
			Set(ctx, c.ns+key, b, time.Duration(expiration)*time.Second).
			Err()
	}
}

func (c *CacheImpl) Get(ctx context.Context, key string) ([]byte, error) {
	b, err := c.client.Get(ctx, c.ns+key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}

		return nil, err
	}
	return b, nil
}

func (c *CacheImpl) GetObject(ctx context.Context, key string, doc interface{}) error {
	b, err := c.client.Get(ctx, c.ns+key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(b, doc)
}

func (c *CacheImpl) GetString(ctx context.Context, key string) (string, error) {
	s, err := c.client.Get(ctx, c.ns+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}
		return "", err
	}
	return s, nil
}

func (c *CacheImpl) GetInt(ctx context.Context, key string) (int64, error) {
	i, err := c.client.Get(ctx, c.ns+key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return i, nil
}

func (c *CacheImpl) GetFloat(ctx context.Context, key string) (float64, error) {
	f, err := c.client.Get(ctx, c.ns+key).Float64()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return f, nil
}

func (c *CacheImpl) Exist(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, c.ns+key).
		Val() > 0
}

func (c *CacheImpl) Delete(ctx context.Context, key string, opts ...DeleteOptions) error {
	deleteCache := &DeleteCache{}
	for _, opt := range opts {
		opt(deleteCache)
	}

	if deleteCache.Pattern != "" {
		iter := c.client.Scan(ctx, 0, c.ns+deleteCache.Pattern, 0).Iterator()

		var localKeys []string
		for iter.Next(ctx) {
			localKeys = append(localKeys, iter.Val())
		}

		if err := iter.Err(); err != nil {
			return err
		}

		if len(localKeys) > 0 {
			_, err := c.client.Pipelined(ctx, func(p redis.Pipeliner) error {
				p.Del(ctx, localKeys...)
				return nil
			})

			if err != nil {
				return err
			}
		}

		return nil
	}

	return c.client.Del(ctx, c.ns+key).Err()
}

func (c *CacheImpl) GetKeys(ctx context.Context, pattern string) []string {
	cmd := c.client.Keys(ctx, pattern)
	keys, err := cmd.Result()
	if err != nil {
		return nil
	}

	return keys
}

func (c *CacheImpl) RemainingTime(ctx context.Context, key string) int {
	return int(c.client.TTL(ctx, c.ns+key).Val().Seconds())
}

func (c *CacheImpl) Close() error {
	return c.client.Close()
}
