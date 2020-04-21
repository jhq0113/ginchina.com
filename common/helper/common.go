package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type entry struct {
	lock     *uint32
	createAt int64
}

func newEntry(createAt int64) *entry {
	var lock uint32

	return &entry{
		lock:     &lock,
		createAt: createAt,
	}
}

var (
	atomicMap sync.Map
)

//region 1.0 Convert
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

func ToInt(value interface{}) int64 {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return val.Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(val.Uint())
	case float32, float64:
		return int64(val.Float())
	case string:
		res, _ := strconv.ParseInt(val.String(), 0, 64)
		return res
	}

	return 0
}

func ToUint(value interface{}) uint64 {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return uint64(val.Int())
	case uint, uint8, uint16, uint32, uint64:
		return val.Uint()
	case float32, float64:
		return uint64(val.Float())
	case string:
		res, _ := strconv.ParseInt(val.String(), 0, 64)
		return uint64(res)
	}

	return 0
}

func Int64ToSecond(val int64) time.Duration {
	return time.Duration(val) * time.Second
}

func ToFloat(str string) float64 {
	res, _ := strconv.ParseFloat(str, 64)
	return res
}

func ToBool(str string) bool {
	res, _ := strconv.ParseBool(str)
	return res
}

/**
 * ip转换为long值
 */
func Ip2long(ip net.IP) uint32 {
	a := uint32(ip[12])
	b := uint32(ip[13])
	c := uint32(ip[14])
	d := uint32(ip[15])
	return uint32(a<<24 | b<<16 | c<<8 | d)
}

/**
 * long值转换为ip
 */
func Long2ip(ip uint32) net.IP {
	a := byte((ip >> 24) & 0xFF)
	b := byte((ip >> 16) & 0xFF)
	c := byte((ip >> 8) & 0xFF)
	d := byte(ip & 0xFF)
	return net.IPv4(a, b, c, d)
}

//endregion

//region 1.1 Math
func Ceil(numerator, denominator int64) int64 {
	return int64(math.Ceil(float64(numerator) / float64(denominator)))
}

//endregion

func GetTraceId(ctx *gin.Context) int64 {
	return ctx.GetInt64("trace-id")
}

//region 1.3 原子操作
func Acquire(key string, timeout int64) bool {
	ent, _ := atomicMap.LoadOrStore(key, newEntry(time.Now().Unix()))
	value := ent.(*entry)
	if atomic.AddUint32(value.lock, 1) == 1 {
		return true
	}

	if time.Now().Unix()-value.createAt > timeout {
		old := value.lock
		atomic.AddUint32(value.lock, 1)
		if atomic.CompareAndSwapUint32(value.lock, *old, *old+1) {
			ent := newEntry(time.Now().Unix())
			atomicMap.Store(key, ent)
			return atomic.AddUint32(ent.lock, 1) == 1
		}
	}
	return false
}

func Release(key string) {
	atomicMap.Store(key, newEntry(0))
}

//endregion

func RenderHtml(ctx *gin.Context, code int, body []byte) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(code, string(body))
}
