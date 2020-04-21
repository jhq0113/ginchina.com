package common

import "encoding/json"

const (
	config = `{
  "application": {
    "addr": ":80",
    "readTimeout": 20,
    "writeTimeout": 10,
    "maxHeaderBytes": 1048576
  },
  "cache": {
    "type":"permap",
    "permap": {
      "fileName": "/tmp/permap.bin",
      "checkInterval": 10
    }
  },
  "log": {
    "fileName": "/tmp/logs/store.log",
    "level": -1
  },
  "params" : {
    "minPageSize": 2,
    "maxPageSize": 1000
  }
}`
)

var (
	DefaultConfig Config
)

func init() {
	_ = json.Unmarshal([]byte(config), &DefaultConfig)
}

/**
 * 应用配置
 */
type Application struct {
	Addr           string `json:"addr"`
	GinMode        string `json:"ginMode"`
	MaxCpu         int    `json:"maxCpu"`
	ReadTimeout    int64  `json:"readTimeout"`
	WriteTimeout   int64  `json:"writeTimeout"`
	MaxHeaderBytes int    `json:"maxHeaderBytes"`
}

/**
 * 数据库配置
 */
type Db struct {
	Driver        string `json:"driver"`
	Dsn           string `json:"dsn"`
	MaxConnection int    `json:"maxConnection"`
	MaxIdle       int    `json:"maxIdle"`
	MaxLifetime   int64  `json:"maxLifetime"`
}

/**
 * 缓存配置
 */
type Cache struct {
	Type   string      `json:"type"`
	Redis  Redis       `json:"redis"`
	Permap PermapCache `json:"permap"`
}

/**
 * 文件缓存配置
 */
type PermapCache struct {
	FileName      string `json:"fileName"`
	CheckInterval int64  `json:"checkInterval"`
}

/**
 * Redis连接配置
 */
type Redis struct {
	Host         string `json:"host"`
	Port         uint16 `json:"port"`
	Db           int    `json:"db"`
	Password     string `json:"password"`
	MaxRetries   int    `json:"maxRetries"`
	PoolSize     int    `json:"poolSize"`
	MinIdleConns int    `json:"minIdleConns"`
	MaxConnAge   int64  `json:"maxConnAge"`
}

/**
 * 日志配置
 */
type Log struct {
	FileName string `json:"fileName"` //日志文件
	Level    int8   `json:"level"`    //日志级别，-1 Debug, 0 Info, 1 Warn, 2 Error, 3 Panic, 4 Fatal
	FileMode uint32 `json:"fileMode"`
}

/**
 * 项目配置
 */
type Config struct {
	Application Application            `json:"application"`
	Cache       Cache                  `json:"cache"`
	Db          Db                     `json:"db"`
	Log         Log                    `json:"log"`
	Params      map[string]interface{} `json:"params"`
}
