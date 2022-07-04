package cmd

import (
	"context"
	"github.com/cloudSlit/cloudslit/ca/api"
	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// InitHTTPServer Initialize HTTP service
func InitHTTPServer(ctx context.Context, handler http.Handler) func() {
	addr := core.Is.Config.HTTP.Listen
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Infof("HTTP server is running at %s.", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(30))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(err.Error())
		}
	}
}

// Run Running services
func RunHttp(ctx context.Context) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	app := api.Serve()
	cleanFunc := InitHTTPServer(ctx, app)

EXIT:
	for {
		sig := <-sc
		logger.Infof("Received signal[%s]", sig.String())
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
	logger.Infof("HTTP service exit")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
