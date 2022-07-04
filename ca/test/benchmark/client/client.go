package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/cloudSlit/cloudslit/ca/pkg/caclient"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/cloudSlit/cloudslit/ca/pkg/spiffe"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	mode     = flag.String("mode", "http", "")
	ca       = flag.String("ca", "", "")
	ocspAddr = flag.String("ocsp", "", "")
	addr     = flag.String("addr", "0.0.0.0:8083", "")
	server   = flag.String("server", "127.0.0.1", "")
	normal   = flag.Bool("normal", false, "")
)

func main() {
	flag.Parse()

	logger.GlobalConfig(logger.Conf{
		Debug: true,
		Level: zapcore.DebugLevel,
	})

	go func() {
		http.ListenAndServe("0.0.0.0:8360", nil)
	}()

	var normalClient *http.Client
	client := &fasthttp.Client{}

	normalTr := &http.Transport{
		TLSHandshakeTimeout: 5 * time.Second,
		MaxIdleConns:        2000,
		MaxIdleConnsPerHost: 200,
		MaxConnsPerHost:     2000,
		IdleConnTimeout:     20 * time.Second,
		WriteBufferSize:     1024,
		ReadBufferSize:      1024,
	}

	switch *mode {
	case "http":
	case "https":
		mtlsConfig := mTlsConfig()
		tlsConfig := &tls.Config{
			RootCAs: mtlsConfig.RootCAs,
		}
		client.TLSConfig = tlsConfig
		normalTr.TLSClientConfig = tlsConfig
	case "mtls":
		mtlsConfig := mTlsConfig()
		client.TLSConfig = mtlsConfig
		normalTr.TLSClientConfig = mtlsConfig
	default:
		panic("unknown mode")
	}

	if client.TLSConfig != nil {
		client.TLSConfig.ClientSessionCache = tls.NewLRUClientSessionCache(64)
	}

	normalClient = &http.Client{
		Transport: normalTr,
	}

	srv := &fasthttp.Server{
		IdleTimeout:        30 * time.Second,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 30 * time.Second,
		MaxConnsPerIP:      200,
	}
	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		scheme := "http"
		if *mode == "https" || *mode == "mtls" {
			scheme = "https"
		}

		if *normal {
			resp, err := normalClient.Get(scheme + "://" + *server)
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			ctx.SetStatusCode(resp.StatusCode)
			ctx.SetBody(body)
			return
		}

		req := &ctx.Request
		req.URI().SetHost(*server)

		req.URI().SetScheme(scheme)
		err := client.Do(req, &ctx.Response)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("started.")

	if err := srv.ListenAndServe(*addr); err != nil {
		panic(err)
	}

	select {}
}

func mTlsConfig() *tls.Config {
	cai := caclient.NewCAI(
		caclient.WithCAServer(caclient.RoleDefault, *ca),
		caclient.WithOcspAddr(*ocspAddr))
	ex, err := cai.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "benchmark_client",
	})
	if err != nil {
		panic(err)
	}
	_, err = ex.Transport.GetCertificate()
	if err != nil {
		panic(err)
	}
	tlsConfig, err := ex.ClientTLSConfig("")
	if err != nil {
		panic(err)
	}
	return tlsConfig.TLSConfig()
}
