package logger

import (
	"fmt"
	"github.com/flowshield/flowshield/ca/pkg/logger/redis_hook"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type ciCore struct {
	zapcore.LevelEnabler

	RedisHook *redis_hook.RedisHook

	fields map[string]interface{}
}

var _ zapcore.Core = (*ciCore)(nil)

// NewCiCore Create custom CiCore
func NewCiCore(hook *redis_hook.RedisHook) zapcore.Core {
	return newCiCore(hook)
}

func newCiCore(hook *redis_hook.RedisHook) *ciCore {
	core := &ciCore{
		LevelEnabler: zapcore.InfoLevel,
		RedisHook:    hook,
		fields:       make(map[string]interface{}),
	}
	return core
}

func (c *ciCore) getAllFields() map[string]interface{} {
	return c.fields
}

func (c *ciCore) Enabled(lvl zapcore.Level) bool {
	if lvl < zapcore.InfoLevel {
		return false
	}
	return true
}

func (c *ciCore) combineFields(fields []zapcore.Field) map[string]interface{} {
	// Copy our map.
	m := make(map[string]interface{}, len(c.fields)+len(fields))
	for k, v := range c.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fields {
		f.AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return m
}

func (c *ciCore) With(fields []zapcore.Field) zapcore.Core {
	m := c.combineFields(fields)
	return &ciCore{
		LevelEnabler: c.LevelEnabler,
		RedisHook:    c.RedisHook,
		fields:       m,
	}
}

func (c *ciCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

func (c *ciCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	var callerStack strings.Builder
	if e.Caller.Defined {
		callerStack.WriteString(e.Caller.TrimmedPath())
		callerStack.WriteString(" ")
		callerStack.WriteString(e.Caller.Function)
	}
	if e.Stack != "" {
		callerStack.WriteString("\n")
		callerStack.WriteString(e.Stack)
	}
	fields = append(fields,
		zap.String("stack", callerStack.String()),
	)

	loggerName := "default"
	if e.LoggerName != "" {
		loggerName = e.LoggerName
	}

	var message strings.Builder
	message.WriteString(loggerName)
	message.WriteString(": ")
	message.WriteString(e.Message)

	combinedFields := c.combineFields(fields)

	additionFields := make(map[string]interface{}, len(combinedFields))
	for field, value := range combinedFields {
		additionFields[field] = value
	}

	// Optimize extra output form [extra1]: data, [extra2]: data2
	if len(additionFields) > 0 {
		for k, v := range additionFields {
			var output string
			if str, ok := v.(string); ok {
				output = str
			} else {
				output, _ = jsoniter.MarshalToString(v)
			}
			message.WriteString(" [")
			message.WriteString(k)
			message.WriteString("]: ")
			message.WriteString(output)
		}
	}

	(&e).Message = message.String()

	msg := redis_hook.CreateZapOriginLogMessage(&e, combinedFields)

	// Marshal into json message
	js, err := jsoniter.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error creating message for REDIS: %s", err)
	}

	// get connection from pool
	conn := c.RedisHook.RedisPool.Get()
	defer conn.Close()

	// send message
	_, err = conn.Do("RPUSH", c.RedisHook.RedisKey, js)
	if err != nil {
		fmt.Fprintln(os.Stdout, "stash log: ", string(js))
		return fmt.Errorf("error sending message to REDIS: %s", err)
	}

	if c.RedisHook.TTL != 0 {
		_, err = conn.Do("EXPIRE", c.RedisHook.RedisKey, c.RedisHook.TTL)
		if err != nil {
			return fmt.Errorf("error setting TTL to key: %s, %s", c.RedisHook.RedisKey, err)
		}
	}

	return nil
}

func (c *ciCore) Sync() error {
	return nil
}
