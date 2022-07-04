package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/cloudSlit/cloudslit/casdk/caclient"
	"github.com/cloudSlit/cloudslit/casdk/pkg/logger"
	"github.com/cloudSlit/cloudslit/casdk/pkg/spiffe"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	caAddr     = flag.String("ca", "https://192.168.2.80:8681", "CA Server")
	ocspAddr   = flag.String("ocsp", "http://192.168.2.80:8682", "Ocsp Server")
	serverAddr = flag.String("server", "https://127.0.0.1:6066", "")
	authKey    = "0739a645a7d6601d9d45f6b237c4edeadad904f2fce53625dfdd541ec4fc8134"
)

// go run server.go -ca https://127.0.0.1:8081 -ocsp http://127.0.0.1:8082 -server https://127.0.0.1:6066

func init() {
	logger.GlobalConfig(logger.Conf{
		Debug: true,
		Level: zapcore.DebugLevel,
	})
}

func main() {
	flag.Parse()
	client, err := NewMTLSClient()
	if err != nil {
		logger.Fatalf("Client init error: %v", err)
	}
	ticker := time.Tick(time.Second)
	for i := 0; i < 1000; i++ {
		<-ticker

		resp, err := client.Get(*serverAddr)
		if err != nil {
			logger.With("resp", resp).Error(err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Infof("Request result: %v, %s", resp.StatusCode, body)
	}
}

// mTLS Client Use example
func NewMTLSClient() (*http.Client, error) {
	l, _ := logger.NewZapLogger(&logger.Conf{
		// Level: 2,
		Level: -1,
	})
	c := caclient.NewCAI(
		caclient.WithCAServer(caclient.RoleDefault, *caAddr),
		caclient.WithAuthKey(authKey),
		caclient.WithOcspAddr(*ocspAddr),
		caclient.WithLogger(l),
	)
	ex, err := c.NewExchanger(&spiffe.IDGIdentity{
		SiteID:    "test_site",
		ClusterID: "cluster_test",
		UniqueID:  "client1",
	})
	if err != nil {
		return nil, errors.Wrap(err, "Exchanger initialization failed")
	}
	cfger, err := ex.ClientTLSConfig("supreme")
	if err != nil {
		panic(err)
	}
	cfger.BindExtraValidator(func(identity *spiffe.IDGIdentity) error {
		fmt.Println("id: ", identity.String())
		return nil
	})
	tlsCfg := cfger.TLSConfig()
	//tlsCfg.VerifyConnection = func(state tls.ConnectionState) error {
	//	cert := state.PeerCertificates[0]
	//	fmt.Println("Server certificate generation time: ", cert.NotBefore.String())
	//	return nil
	//}
	client := httpClient(tlsCfg)
	go ex.RotateController().Run()
	// util.ExtractCertFromExchanger(ex)

	resp, err := client.Get("http://www.google.com")
	if err != nil {
		panic(err)
	}

	fmt.Println("baidu test: ", resp.StatusCode)

	return client, nil
}

func httpClient(cfg *tls.Config) *http.Client {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig:   cfg,
			DisableKeepAlives: true,
		},
	}
	return &client
}
