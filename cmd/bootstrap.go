package cmd

import (
	"github/yogabagas/join-app/config"
	"github/yogabagas/join-app/pkg/cache/redis"
	"github/yogabagas/join-app/pkg/database/sql"
	"github/yogabagas/join-app/shared/constant"
	"net/url"
)

var (
	configURL string

	sqlDB       = sql.DBConn
	redisClient = redis.CacheConn
)

func InitSQLModule() (*sql.DB, error) {
	return sql.NewDBConn(constant.MySQL.String())
}

func InitCache() (*redis.Cache, error) {
	redisCreds := url.URL{
		Host: config.GlobalCfg.Cache.Redis.Host,
		User: url.UserPassword(config.GlobalCfg.Cache.Redis.User, config.GlobalCfg.Cache.Redis.Password),
	}

	return redis.NewCache(&redisCreds)
}
