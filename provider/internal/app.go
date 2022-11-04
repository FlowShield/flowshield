package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/server"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/mysql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	if err = mysql.Init(); err != nil {
		return nil, err
	}
	if err = server.InitNode(ctx); err != nil {
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
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()
	ps := http.NewServeMux()
	ps.Handle("/metrics", promhttp.Handler())
	ps.HandleFunc("/debug/pprof/", pprof.Index)
	ps.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	ps.HandleFunc("/debug/pprof/profile", pprof.Profile)
	ps.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	ps.HandleFunc("/debug/pprof/trace", pprof.Trace)
	pSrv := &http.Server{
		Addr:        config.C.App.HttpListenAddr,
		Handler:     ps,
		ReadTimeout: 5 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	go func() {
		logger.Infof("pprof server is running at %s.", config.C.App.HttpListenAddr)
		err := pSrv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
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
