package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/cloudslit/cloudslit/provider/internal/bll"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/internal/initer"
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/influxdb"
	influx_client "github.com/cloudslit/cloudslit/provider/pkg/influxdb/client/v2"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/cloudslit/cloudslit/provider/pkg/util/structure"
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
	err = initer.InitMachine()
	if err != nil {
		return nil, err
	}
	// initialize the timing module
	influxdbCleanFunc, err := InitInfluxdb(ctx)
	if err != nil {
		return nil, err
	}
	InitHttpClient()
	return func() {
		loggerCleanFunc()
		influxdbCleanFunc()
	}, nil
}

// 启动服务
func InitServer(ctx context.Context, opts ...Option) (func(), error) {
	initCleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return nil, err
	}
	basicConf, attr, err := initer.InitCert([]byte(config.C.Certificate.CertPem))
	if err != nil {
		return nil, err
	}
	switch basicConf.Type {
	//case initer.TypeClient:
	//	//serverCleanFunc = bll.NewClient().Listen(ctx, attr)
	//	fmt.Println("########## start the client proxy #########")
	//case initer.TypeServer:
	//	//bll.NewServer().Listen(ctx, attr)
	//	fmt.Println("########## start the server proxy #########")
	case initer.TypeRelay:
		fmt.Println("########## start the relay proxy #########")
		bll.NewRelay().Listen(ctx, attr)
	default:
		return nil, errors.New("error type")
	}
	return func() {
		initCleanFunc()
	}, nil
}

func InitInfluxdb(ctx context.Context) (func(), error) {
	if !config.C.Influxdb.Enabled {
		logger.WithContext(ctx).Warn("Influxdb Function is disabled")
		return func() {}, nil
	}
	client, err := influx_client.NewHTTPClient(influx_client.HTTPConfig{
		Addr:                fmt.Sprintf("http://%v:%v", config.C.Influxdb.Address, config.C.Influxdb.Port),
		Username:            config.C.Influxdb.Username,
		Password:            config.C.Influxdb.Password,
		MaxIdleConns:        config.C.Influxdb.MaxIdleConns,
		MaxIdleConnsPerHost: config.C.Influxdb.MaxIdleConns,
	})
	if err != nil {
		return func() {}, err
	}
	if _, _, err := client.Ping(1 * time.Second); err != nil {
		_ = client.Close()
		return func() {}, err
	}
	iconfig := new(influxdb.CustomConfig)
	structure.Copy(config.C.Influxdb, iconfig)
	metrics, err := influxdb.NewMetrics(&influxdb.HTTPClient{
		Client: client,
		BatchPointsConfig: influx_client.BatchPointsConfig{
			Precision: config.C.Influxdb.Precision,
			Database:  config.C.Influxdb.Database,
		},
	}, iconfig)
	config.Is.Metrics = metrics
	return func() {
		client.Close()
	}, err
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

func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := InitServer(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		logger.WithContext(ctx).Infof("received signal[%s]", sig.String())
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
	logger.WithContext(ctx).Infof("shutdown!")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
