//go:build linux || darwin || windows
// +build linux darwin windows

package internal

import "C"
import (
	"context"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/initer"
	"github.com/cloudslit/cloudslit/provider/internal/server"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	err := config.MustLoad(o.ConfigFile)
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
	err = server.InitNode(ctx)
	if err != nil {
		return nil, err
	}
	InitHTTPServer(ctx)
	return func() {
		loggerCleanFunc()
	}, nil
}

// InitHTTPServer 初始化http服务
func InitHTTPServer(ctx context.Context) {
	addr := "0.0.0.0:" + strconv.Itoa(config.C.App.LocalPort)
	logger.Infof("Node server is running at %s.", addr)
	http.HandleFunc("/", IndexHandler)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("error: ", err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
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
