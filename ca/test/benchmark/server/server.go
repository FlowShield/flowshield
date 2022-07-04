package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/cloudslit/cloudslit/ca/pkg/caclient"
	"github.com/cloudslit/cloudslit/ca/pkg/spiffe"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

var (
	mode      = flag.String("mode", "http", "")
	ca        = flag.String("ca", "", "")
	ocspAddr  = flag.String("ocsp", "", "")
	addr      = flag.String("addr", "0.0.0.0:28081", "")
	keepAlive = flag.Bool("keepalive", false, "")

	cert *tls.Certificate
)

func main() {
	flag.Parse()

	srv := &fasthttp.Server{
		IdleTimeout:        30 * time.Second,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 30 * time.Second,
		MaxConnsPerIP:      200,
		DisableKeepalive:   *keepAlive,
	}
	var tlsCfg *tls.Config

	switch *mode {
	case "http":
	case "https":
		_ = mTlsConfig()
		tlsCfg = &tls.Config{
			ClientAuth:   tls.NoClientCert,
			Certificates: []tls.Certificate{*cert},
		}
	case "mtls":
		tlsCfg = mTlsConfig()
	default:
		panic("unknown mode")
	}

	ln, err := net.Listen("tcp4", *addr)
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	lnTls := tls.NewListener(ln, tlsCfg)

	srv.Handler = func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte("hello"))
	}

	fmt.Println("started.")

	if tlsCfg != nil {
		if err := srv.Serve(lnTls); err != nil {
			panic(err)
		}
	} else {
		if err := srv.Serve(ln); err != nil {
			panic(err)
		}
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
		UniqueID:  "benchmark_server",
	})
	if err != nil {
		panic(err)
	}
	cert, err = ex.Transport.GetCertificate()
	if err != nil {
		panic(err)
	}
	tlsConfig, err := ex.ServerTLSConfig()
	if err != nil {
		panic(err)
	}
	return tlsConfig.TLSConfig()
}
