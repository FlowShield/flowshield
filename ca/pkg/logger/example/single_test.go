package example

import (
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/flowshield/flowshield/ca/pkg/logger/redis_hook"
	"go.uber.org/zap/zapcore"
	"log"
)

var (
	EnvEnableRedisOutput bool // Simulated environment variables
	EnvDebug             bool
)

func init() {
	EnvEnableRedisOutput = true
	EnvDebug = true
	initLogger()
}

func initLogger() {
	conf := &logger.Conf{
		Level:  zapcore.DebugLevel, // Output log level
		Caller: true,               //Whether to open record calling folder + number of lines + function name
		Debug:  true,               // Enable debug
		// All logs output to redis are above info level
		AppInfo: &logger.ConfigAppData{
			AppVersion: "1.0",
			Language:   "zh-cn",
		},
	}
	if !EnvDebug || EnvEnableRedisOutput {
		// In case of production environment
		conf.Level = zapcore.InfoLevel
		conf.HookConfig = &redis_hook.HookConfig{
			Key:  "log_key",
			Host: "redis.msp",
			Port: 6380,
		}
	}
	err := logger.GlobalConfig(*conf)
	if err != nil {
		log.Print("[ERR] Logger init error: ", err)
	}
	logger.Infof("info test: %v", "data")
}
