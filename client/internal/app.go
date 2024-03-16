package internal

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/flowshield/flowshield/client/internal/config"
	"github.com/flowshield/flowshield/client/pkg/logger"
	"github.com/flowshield/flowshield/client/pkg/web3/w3s"
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
	err := config.MustLoad(o.ConfigFile)
	if err != nil {
		return nil, err
	}
	config.PrintWithJSON()
	logger.WithContext(ctx).Printf("Service started, running mode：%s，version：%s，process number：%d", config.C.RunMode, o.Version, os.Getpid())

	// initialize the log module
	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}
	err = InitMachine()
	if err != nil {
		return nil, err
	}
	if err := w3s.Init(&config.C.Web3); err != nil {
		return nil, err
	}
	// initialize Http Client
	InitHttpClient()
	// start client server
	InitClientServer(ctx)

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

// Run service
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
		logger.WithContext(ctx).Infof("signal received[%s]", sig.String())
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
	logger.WithContext(ctx).Infof("Service exit")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
