package cache

import (
	"encoding/json"
	"github.com/jhq0113/ginchina.com/common"
	"github.com/jhq0113/ginchina.com/common/helper"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"
)

type Permap struct {
	ICache
	data           sync.Map
	lastSyncTime   int64
	newUpdateTimes uint64
}

type item struct {
	Value    []byte `json:"value"`
	ExpireAt int64  `json:"expireAt"`
}

func NewPermap(config *common.PermapCache) *Permap {
	cache := &Permap{}
	if config.CheckInterval < 1 {
		config.CheckInterval = 1
	}

	//从文件中加载缓存
	cache.load(config.FileName)

	go cache.persist(config)

	return cache
}

func (this *Permap) load(fileName string) {
	if !helper.FileExists(fileName) {
		return
	}

	file, _ := helper.FOpen(fileName, "r", 0755)
	defer func() {
		file.Close()
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil || len(data) < 1 {
		return
	}

	persistData := map[string]item{}
	err = json.Unmarshal(data, &persistData)
	if err != nil {
		return
	}

	for key, current := range persistData {
		this.data.Store(key, current)
	}

	this.lastSyncTime = time.Now().Unix()
}

func (this *Permap) persist(config *global.PermapCache) {
	persistTick := time.NewTicker(helper.Int64ToSecond(config.CheckInterval))
	for {
		select {
		case <-persistTick.C:
			if this.newUpdateTimes < 1 && time.Now().Unix()-this.lastSyncTime < 6000 {
				continue
			}

			//当缓存发生更新，或者缓存未发生更新但时间已经过了6000秒
			persistData := map[string]item{}
			this.data.Range(func(key, value interface{}) bool {
				if current, ok := value.(item); ok {
					persistData[key.(string)] = current
				}
				return true
			})

			if len(persistData) < 1 {
				continue
			}

			data, err := json.Marshal(persistData)
			if err != nil {
				continue
			}
			file, err := helper.FOpen(config.FileName, "w+", 0775)
			if err != nil {
				continue
			}
			length, err := file.Write(data)
			if err != nil || length != len(data) {
				continue
			}

			_ = file.Close()
			this.newUpdateTimes = 0
			this.lastSyncTime = time.Now().Unix()
		}
	}
}

func (this *Permap) Get(key string) []byte {
	data, ok := this.data.Load(key)
	if ok {
		current := data.(item)
		if current.ExpireAt < 1 || current.ExpireAt > time.Now().Unix() {
			return current.Value
		}
	}
	return nil
}

func (this *Permap) Set(key string, value []byte, timeout int64) bool {
	current := item{
		Value: value,
	}
	if timeout > 0 {
		current.ExpireAt = time.Now().Unix() + timeout
	}
	this.data.Store(key, current)

	atomic.AddUint64(&this.newUpdateTimes, 1)

	return true
}

func (this *Permap) Exists(key string) bool {
	data, ok := this.data.Load(key)
	if !ok {
		return false
	}
	current := data.(*item)
	return current.ExpireAt < 1 || current.ExpireAt > time.Now().Unix()
}

func (this *Permap) Delete(keys ...string) bool {
	for _, key := range keys {
		this.data.Delete(key)
	}

	atomic.AddUint64(&this.newUpdateTimes, 1)
	return true
}
