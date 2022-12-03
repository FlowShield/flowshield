package gin

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/flowshield/flowshield/verifier/pkg/confer"
	"github.com/gin-gonic/gin"
)

type OnShutdownF struct {
	f       func(cancel context.CancelFunc)
	timeout time.Duration
}

var (
	onShutdown []OnShutdownF
)

func RegisterOnShutdown(f func(cancel context.CancelFunc), timeout time.Duration) {
	onShutdown = append(onShutdown, OnShutdownF{
		f:       f,
		timeout: timeout,
	})
}

func NewGin() *gin.Engine {
	if confer.ConfigEnvIsDev() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	if confer.ConfigEnvIsDev() {
		r.Use(gin.Logger())
	}
	return r
}

func ListenHTTP(httpPort string, r http.Handler, timeout int, f ...func()) {
	srv := &http.Server{
		Addr:    httpPort,
		Handler: r,
	}
	// 监听端口
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// 注册关闭使用函数
	for _, v := range f {
		srv.RegisterOnShutdown(v)
	}
	// 监听信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit
	// 执行on shutdown 函数 - 同步
	for _, v := range onShutdown {
		var wg sync.WaitGroup
		wg.Add(1)
		ctx, cancel := context.WithCancel(context.TODO())
		go v.f(cancel)
		select {
		case <-time.After(v.timeout):
			log.Println("on shutdown timeout:", f)
			wg.Done()
		case <-ctx.Done():
			wg.Done()
		}
		wg.Wait()
	}
	// 执行shutdown
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panic("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
