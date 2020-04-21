package redis

import (
	"github.com/go-redis/redis/v7"
	"github.com/jhq0113/ginchina.com/common"
	"github.com/jhq0113/ginchina.com/common/helper"
)

func NewRedis(options *common.Redis) *redis.Client {
	host := options.Host
	if host == "" {
		host = "127.0.0.1"
	}

	port := options.Port
	if port == 0 {
		port = 6379
	}

	return redis.NewClient(&redis.Options{
		Addr:         helper.Append(host, ":", helper.ToString(port)),
		Password:     options.Password,
		DB:           options.Db,
		MaxRetries:   options.MaxRetries,
		PoolSize:     options.PoolSize,
		MinIdleConns: options.MinIdleConns,
		MaxConnAge:   helper.Int64ToSecond(options.MaxConnAge),
	})
}
