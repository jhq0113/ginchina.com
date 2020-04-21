package cache

import (
	redis2 "github.com/go-redis/redis/v7"
	"github.com/jhq0113/ginchina.com/common/helper"
)

type Redis struct {
	ICache
	client *redis2.Client
}

func NewRedis(cli *redis2.Client) *Redis {
	return &Redis{
		client: cli,
	}
}

func (this *Redis) Get(key string) []byte {
	data, _ := this.client.Get(key).Bytes()
	return data
}

func (this *Redis) Set(key string, value []byte, timeout int64) bool {
	return this.client.Set(key, string(value), helper.Int64ToSecond(timeout)).Err() == nil
}

func (this *Redis) Exists(key string) bool {
	return this.client.Exists(key).Val() == 1
}

func (this *Redis) Delete(keys ...string) bool {
	return this.client.Del(keys...).Err() == nil
}
