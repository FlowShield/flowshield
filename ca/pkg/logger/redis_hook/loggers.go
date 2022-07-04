package redis_hook

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

// zap need extra data for fields
func CreateZapOriginLogMessage(entry *zapcore.Entry, data map[string]interface{}) map[string]interface{} {
	fields := make(map[string]interface{}, len(data))
	if data != nil {
		for k, v := range data {
			fields[k] = v
		}
	}
	var level = strings.ToUpper(entry.Level.String())
	if level == "ERROR" {
		level = "ERR"
	}
	if level == "WARN" {
		level = "WARNING"
	}
	if level == "FATAL" {
		level = "CRIT"
	}
	fields["level"] = level
	fields["message"] = entry.Message
	return fields
}
