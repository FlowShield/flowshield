package internal

import (
	"context"
	"crypto/tls"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/initer"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type options struct {
	ConfigFile string
	ModelFile  string
	Version    string
}

// Option Defining configuration items
type Option func(*options)

// SetConfigFile setting the configuration file
func SetConfigFile(s string) Option {
	return func(o *options) {
		o.ConfigFile = s
	}
}

// SetVersion set version number
func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init application initialization
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.ConfigFile)
	// working with environment variables
	err := config.ParseConfigByEnv()
	if err != nil {
		return nil, err
	}
	config.PrintWithJSON()
	logger.WithContext(ctx).Printf("Service started, running mode：%s，version：%s，process number：%d", config.C.RunMode, o.Version, os.Getpid())

	// initialize the log module
	loggerCleanFunc, err := initer.InitLogger()
	if err != nil {
		return nil, err
	}

	InitHttpClient()
	InitProviderServer(ctx)
	return func() {
		loggerCleanFunc()
	}, nil
}

func InitHttpClient() {
	config.Is.HttpClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			IdleConnTimeout: 5 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

// Run 运行服务
func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("接收到信号[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}

	cleanFunc()
	logger.WithContext(ctx).Infof("服务退出")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
