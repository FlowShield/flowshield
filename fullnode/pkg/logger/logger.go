package logger

import (
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var simpleLogger *zap.SugaredLogger

type Config struct {
	Level       string
	SendToFile  bool
	Filename    string
	NoCaller    bool
	NoLogLevel  bool
	Development bool
	MaxSize     int // megabytes
	MaxAge      int // days
	MaxBackups  int
}

func Init(cfg *Config) {
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(err)
	}

	consoleSyncer := zapcore.AddSync(os.Stdout)
	consoleEncoder := getConsoleEncoder(cfg)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, l)

	var opts []zap.Option
	opts = append(opts, zap.AddStacktrace(zap.DPanicLevel), zap.AddCallerSkip(1))
	if !cfg.NoCaller {
		opts = append(opts, zap.AddCaller())
	}
	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	core := consoleCore
	if cfg.SendToFile {
		fileSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
		fileEncoder := getJSONEncoder(cfg)
		fileCore := zapcore.NewCore(fileEncoder, fileSyncer, l)

		core = zapcore.NewTee(consoleCore, fileCore)
	}

	logger = zap.New(core, opts...)
	simpleLogger = logger.Sugar()
}

func getJSONEncoder(cfg *Config) zapcore.Encoder {
	return getEncoder(cfg, true)
}

func getConsoleEncoder(cfg *Config) zapcore.Encoder {
	return getEncoder(cfg, false)
}

func lastNthIndexString(s string, sub string, index int) string {
	r := strings.Split(s, sub)
	if len(r) < index {
		return s
	}
	return strings.Join(r[len(r)-index:], "/")
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(lastNthIndexString(caller.String(), "/", 3))
}

func getEncoder(cfg *Config, jsonFormat bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeCaller = customCallerEncoder

	if cfg.NoLogLevel {
		encoderConfig.LevelKey = zapcore.OmitKey
	}

	if jsonFormat {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func NopSugaredLogger() *zap.SugaredLogger {
	return zap.NewNop().Sugar()
}

func Logger() *zap.Logger {
	return getLogger()
}

func SugaredLogger() *zap.SugaredLogger {
	return getSugaredLogger()
}

func getLogger() *zap.Logger {
	if logger == nil {
		panic("Logger is not initialized yet!")
	}

	return logger
}

func getSugaredLogger() *zap.SugaredLogger {
	return getLogger().Sugar()
}

func getSimpleLogger() *zap.SugaredLogger {
	if simpleLogger == nil {
		panic("Logger is not initialized yet!")
	}
	simpleLogger.With()
	return simpleLogger
}

func NewFileLogger(path string) *zap.Logger {
	fileSyncer := getLogWriter(path, 0, 0, 0)
	fileEncoder := getJSONEncoder(&Config{})
	fileCore := zapcore.NewCore(fileEncoder, fileSyncer, zap.DebugLevel)
	return zap.New(fileCore)
}

func With(fields ...zap.Field) *zap.Logger {
	return getLogger().With(fields...)
}

func Debug(args ...interface{}) {
	getSimpleLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	getSimpleLogger().Debugf(format, args...)
}

func Info(args ...interface{}) {
	getSimpleLogger().Info(args...)
}

func Infof(format string, args ...interface{}) {
	getSimpleLogger().Infof(format, args...)
}

func Warn(c *gin.Context, args ...interface{}) {
	//getSimpleLogger().Warn(args...)
	getLogger().With(ginField(c)...).Sugar().Warn(args...)
}

func Warnf(c *gin.Context, format string, args ...interface{}) {
	//getSimpleLogger().Warnf(format, args...)
	getLogger().With(ginField(c)...).Sugar().Warnf(format, args...)
}

func Error(c *gin.Context, args ...interface{}) {
	getLogger().With(ginField(c)...).Sugar().Error(args...)
	//getSimpleLogger().Error(args...)
}

func Errorf(c *gin.Context, format string, args ...interface{}) {
	getLogger().With(ginField(c)...).Sugar().Errorf(format, args...)
	//getSimpleLogger().Errorf(format, args...)
}

func DPanic(args ...interface{}) {
	getSimpleLogger().DPanic(args...)
}

func DPanicf(format string, args ...interface{}) {
	getSimpleLogger().DPanicf(format, args...)
}

func Panic(args ...interface{}) {
	getSimpleLogger().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	getSimpleLogger().Panicf(format, args...)
}

func Fatal(args ...interface{}) {
	getSimpleLogger().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	getSimpleLogger().Fatalf(format, args...)
}

const loggerKey = iota

func NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(loggerKey), WithContext(ctx).With(fields...))
}

func WithContext(ctx *gin.Context) *zap.Logger {
	if ctx == nil {
		return Logger()
	}
	l, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := l.(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return Logger()
}

func ginField(c *gin.Context) []zap.Field {
	if c != nil {
		return []zap.Field{zap.String("url", c.Request.Method+": "+c.Request.URL.Path)}
	}
	return nil
}
