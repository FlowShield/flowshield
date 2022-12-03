package cmd

import (
	"context"
	"crypto/tls"
	"github.com/flowshield/flowshield/ca/ca/keymanager"
	"github.com/flowshield/flowshield/ca/ca/singleca"
	"github.com/flowshield/flowshield/ca/core"
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// InitTlsServer Initialize TLS service
func InitTlsServer(ctx context.Context, handler *mux.Router) func() {
	addr := core.Is.Config.HTTP.CaListen
	tlsCfg := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return keymanager.GetKeeper().GetCachedTLSKeyPair()
		},
		InsecureSkipVerify: true,
		ClientAuth:         tls.NoClientCert,
	}
	srv := &http.Server{
		Addr:         addr,
		TLSConfig:    tlsCfg,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Infof("TLS server is running at %s.", addr)
		err := srv.ListenAndServeTLS("", "")
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

func RunTls(ctx context.Context) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	app, err := singleca.Server()
	if err != nil {
		return err
	}
	cleanFunc := InitTlsServer(ctx, app)

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
	logger.Infof("TLS service exit")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
