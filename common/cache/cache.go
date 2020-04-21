package cache

import (
	"encoding/json"
	"github.com/jhq0113/ginchina.com/common"
	"github.com/jhq0113/ginchina.com/common/redis"
)

type ICache interface {
	Set(key string, value []byte, timeout int64) bool
	Get(key string) []byte
	Exists(key string) bool
	Delete(keys ...string) bool
}

var cache ICache

const (
	FILE  = "file"
	REDIS = "redis"
)

func Start() {
	switch common.Config.Cache.Type {
	case REDIS:
		cache = NewRedis(redis.NewRedis(&common.Config.Cache.Redis))
	default:
		cache = NewPermap(&common.Config.Cache.Permap)
	}
}

func decode(data []byte, out interface{}) {
	json.Unmarshal(data, out)
}

func encode(data interface{}) []byte {
	encodeData, _ := json.Marshal(data)
	return encodeData
}

func Set(key string, value interface{}, timeout int64) bool {
	return cache.Set(key, encode(value), timeout)
}

func Get(key string, out interface{}) {
	result := cache.Get(key)
	if len(result) < 1 {
		return
	}

	decode(result, out)
}

func Exists(key string) bool {
	return cache.Exists(key)
}

func Delete(keys ...string) bool {
	return cache.Delete(keys...)
}
