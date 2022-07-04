package test

import (
	"crypto/tls"
	"fmt"
	"github.com/cloudSlit/cloudslit/ca/pkg/caclient"
	"github.com/cloudSlit/cloudslit/ca/pkg/spiffe"
	"github.com/valyala/fasthttp"
	"github.com/ztalab/cfssl/helpers"
	cflog "github.com/ztalab/cfssl/log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMTls(t *testing.T) {
	cflog.Level = cflog.LevelDebug
	c := caclient.NewCAI(
		caclient.WithCAServer(caclient.RoleDefault, "https://127.0.0.1:8081"),
		caclient.WithOcspAddr("http://127.0.0.1:8082"))
	serverEx, err := c.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "server1",
	})
	clientEx, err := c.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "client1",
	})
	if err != nil {
		t.Error("transport Error: ", err)
	}

	serverTls, err := serverEx.ServerTLSConfig()
	if err != nil {
		t.Error("Server TLS get error: ", err)
	}
	fmt.Println("------------- Server trust certificate --------------")
	fmt.Println(string(helpers.EncodeCertificatesPEM(serverEx.Transport.ClientTrustStore.Certificates())))
	fmt.Println("------------- END Server trust certificate --------------")

	clientTls, err := clientEx.ClientTLSConfig("")
	if err != nil {
		t.Error("client tls config get error: ", err)
	}
	fmt.Println("------------- Client trust certificate --------------")
	fmt.Println(string(helpers.EncodeCertificatesPEM(clientEx.Transport.TrustStore.Certificates())))
	fmt.Println("------------- END Client trust certificate --------------")

	go func() {
		httpsServer(serverTls.TLSConfig())
	}()
	client := httpClient(clientTls.TLSConfig())
	time.Sleep(2 * time.Second)

	var messages = []string{"hello world", "hello", "world"}
	for range messages {
		resp, err := client.Get("https://127.0.0.1:8082/test111111")
		if err != nil {
			fmt.Fprint(os.Stderr, "request was aborted: ", err)
		}

		fmt.Println("Request succeeded: ", resp.Status)
	}
}

func httpClient(cfg *tls.Config) *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     cfg,
			MaxIdleConns:        50,
			MaxIdleConnsPerHost: 50,
		},
	}
	return &client
}

func httpsServer(cfg *tls.Config) {
	ln, err := net.Listen("tcp4", "0.0.0.0:8082")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	lnTls := tls.NewListener(ln, cfg)

	if err := fasthttp.Serve(lnTls, func(ctx *fasthttp.RequestCtx) {
		str := ctx.Request.String()
		fmt.Println("Server reception: ", str)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"))
	}); err != nil {
		panic(err)
	}
}
