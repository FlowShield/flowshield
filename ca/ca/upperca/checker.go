package upperca

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/ztalab/cfssl/api/client"
	"go.uber.org/zap"

	"github.com/cloudslit/cloudslit/ca/ca/keymanager"
	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/logic/schema"
	"github.com/cloudslit/cloudslit/ca/pkg/influxdb"
)

const CfsslHealthApi = "/api/v1/cfssl/health"

type Checker interface {
	Run()
}

var httpClient = resty.NewWithClient(&http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	},
	Timeout: 5 * time.Second,
})

type checker struct {
	keymanager.UpperClients
	logger *zap.SugaredLogger
	influx *influxdb.Metrics
}

func (hc *checker) Run() {
	if !core.Is.Config.Influxdb.Enabled {
		return
	}
	go func() {
		for {
			<-time.After(5 * time.Second)
			for _, upperClient := range hc.UpperClients.AllClients() {
				go hc.checkUpper(upperClient)
			}
		}
	}()
}

func (hc *checker) checkUpper(upperClient *client.AuthRemote) {
	var err error
	caUrl := upperClient.Hosts()[0]
	caHost := schema.GetHostFromUrl(caUrl)

	resp, err := httpClient.R().Get(caUrl + CfsslHealthApi)
	// Statistical trust certificate comb
	var statusCode int
	if err != nil {
		statusCode = 599
	} else {
		statusCode = resp.StatusCode()
	}

	if err != nil || resp.StatusCode() != http.StatusOK {
		hc.logger.Warnf("Upper CA: %s Connection error: %s", caHost, err)
	}

	hc.influx.AddPoint(&influxdb.MetricsData{
		Measurement: schema.MetricsUpperCaInfo,
		Fields: map[string]interface{}{
			"status": statusCode,
			"delay":  resp.Time().Milliseconds(), // ms
		},
		Tags: map[string]string{
			"host":   caHost,
			"status": strconv.Itoa(statusCode),
			"ip":     schema.GetLocalIpLabel(),
		},
	})
}

// NewChecker Execute only in subordinate CAS
func NewChecker() Checker {
	return &checker{
		UpperClients: keymanager.GetKeeper().RootClient,
		logger:       logger.Named("upper").SugaredLogger,
		influx:       core.Is.Metrics,
	}
}
