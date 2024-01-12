package redis

import (
	"context"
	"crypto/tls"
	"net/url"
	"strings"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
)

const (
	schemaDefault = "redis"
)

type Cache struct {
	Client *redis.Client
	ns     string
}

var CacheConn *Cache

func NewCache(url *url.URL) (*Cache, error) {

	if CacheConn == nil {
		pass, _ := url.User.Password()
		opt := &redis.Options{
			Addr:     url.Host,
			Password: pass,
			DB:       0,
		}

		getTls := url.Query().Get("tls")
		emptyTls := getTls == ""
		if !emptyTls {
			opt.TLSConfig = &tls.Config{ServerName: getTls}
		}

		redisClient := redis.NewClient(opt)
		redisClient.AddHook(redisotel.TracingHook{})

		ns := strings.TrimSuffix(url.Path, "/")
		if ns == "" {
			ns = schemaDefault
		}

		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			return nil, err
		}

		CacheConn = &Cache{
			Client: redisClient,
			ns:     strings.TrimPrefix(url.Path, "/"),
		}
	}

	return CacheConn, nil
}
