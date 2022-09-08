package confer

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var globalConfig *Server
var mutex sync.RWMutex

type Server struct {
	App   App   `mapstructure:"app" json:"app" yaml:"app"`
	Code  Code  `mapstructure:"code" json:"code" yaml:"code"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Log   Log   `mapstructure:"log" json:"log" yaml:"log"`
	sync.RWMutex
}

type App map[string]interface{}

type Code map[string]interface{}

type Redis struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Prefix  string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
}

type Mysql struct {
	Enabled bool     `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	DBName  string   `mapstructure:"dbname" json:"dbName" yaml:"dbname"`
	Prefix  string   `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Pool    DBPool   `mapstructure:"pool" json:"pool" yaml:"pool"`
	Write   DBBase   `mapstructure:"write" json:"write" yaml:"write"`
	Reads   []DBBase `mapstructure:"reads" json:"reads" yaml:"reads"`
}

type DBPool struct {
	PoolMinCap      int   `mapstructure:"pool-min-cap" json:"poolMinCap" yaml:"pool-min-cap"`
	PoolExCap       int   `mapstructure:"pool-ex-cap" json:"poolExCap" yaml:"pool-ex-cap"`
	PoolMaxCap      int   `mapstructure:"pool-max-cap" json:"pool-max-cap" yaml:"pool-max-cap"`
	PoolIdleTimeout int   `mapstructure:"pool-idle-timeout" json:"poolIdleTimeout" yaml:"pool-idle-timeout"`
	PoolWaitCount   int64 `mapstructure:"pool-wait-count" json:"poolWaitCount" yaml:"pool-wait-count"`
	PoolWaitTimeout int   `mapstructure:"pool-wai-timeout" json:"poolWaitTimeout" yaml:"pool-wai-timeout"`
}

type DBBase struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `json:"-"`
	Prefix   string `json:"-"`
}

type Log struct {
	Enabled bool         `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	OutPut  string       `mapstructure:"out-put" json:"outPut" yaml:"out-put"`
	Debug   bool         `mapstructure:"debug" json:"debug" yaml:"debug"`
	Key     string       `mapstructure:"key" json:"key" yaml:"key"`
	Level   logrus.Level `mapstructure:"level" json:"level" yaml:"level"`
	Redis   struct {
		Host string
		Port int
	}
	App struct {
		AppName    string `mapstructure:"app-name" json:"appName" yaml:"app-name"`
		AppID      string `mapstructure:"app-id" json:"appID" yaml:"app-id"`
		AppVersion string `mapstructure:"app-version" json:"appVersion" yaml:"app-version"`
		AppKey     string `mapstructure:"app-key" json:"appKey" yaml:"app-key"`
		Channel    string `mapstructure:"channel" json:"channel" yaml:"channel"`
		SubOrgKey  string `mapstructure:"sub-org-key" json:"subOrgKey" yaml:"sub-org-key"`
		Language   string `mapstructure:"language" json:"language" yaml:"language"`
	} `mapstructure:"app" json:"app" yaml:"app"`
}

type APM struct {
	Addr  string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Http  bool   `mapstructure:"http" json:"http" yaml:"http"`
	Mysql bool   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis bool   `mapstructure:"redis" json:"redis" yaml:"redis"`
}

func Init(configURL string) (err error) {
	f, err := os.Open(configURL)
	if err != nil {
		return
	}
	if err = yaml.NewDecoder(f).Decode(&globalConfig); err != nil {
		return
	}
	if err = handleConfig(globalConfig); err != nil {
		return
	}
	return
}

func handleConfig(config *Server) (err error) {
	config.replaceByEnv(&config.Redis.Address)
	config.replaceByEnv(&config.Mysql.DBName)
	config.replaceByEnv(&config.Mysql.Write.Host)
	// 处理mysql地址
	host, port, err := net.SplitHostPort(globalConfig.Mysql.Write.Host)
	if err != nil {
		err = fmt.Errorf("mysql host port is wrong :%w,%s", err, globalConfig.Mysql.Write.Host)
		return err
	}
	config.Mysql.Write.Host = host
	portInt, _ := strconv.Atoi(port)
	config.Mysql.Write.Port = portInt
	config.replaceByEnv(&config.Mysql.Write.User)
	config.replaceByEnv(&config.Mysql.Write.Password)
	config.Mysql.Write.DBName = config.Mysql.DBName
	config.Mysql.Write.Prefix = config.Mysql.Prefix
	config.replaceByEnv(&config.Log.Redis.Host)
	config.replaceByEnv(&config.Log.App.AppName)
	config.replaceByEnv(&config.Log.App.AppID)
	config.replaceByEnv(&config.Log.App.AppVersion)
	config.replaceByEnv(&config.Log.App.AppKey)
	config.replaceByEnv(&config.Log.App.Channel)
	config.replaceByEnv(&config.Log.App.SubOrgKey)
	config.replaceByEnv(&config.Log.App.Language)
	return
}

func GlobalConfig() *Server {
	mutex.RLock()
	defer mutex.RUnlock()
	return globalConfig
}

func (*Server) replaceByEnv(conf *string) {
	if s := os.Getenv(*conf); len(s) > 0 {
		*conf = s
	}
}
