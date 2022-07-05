package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewLogger(t *testing.T) {
	defer Sync()
	GlobalConfig(Conf{
		Debug:  true,
		Caller: true,
		AppInfo: &ConfigAppData{
			AppName:    "test",
			AppID:      "test",
			AppVersion: "1.0",
			AppKey:     "test",
			Channel:    "1",
			SubOrgKey:  "key",
			Language:   "zh",
		},
	})
	S().Info("test")
}

func TestColorLogger(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := config.Build()

	logger.Info("Now logs should be colored")
}
