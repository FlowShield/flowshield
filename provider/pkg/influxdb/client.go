package influxdb

import (
	"github.com/cloudslit/cloudslit/provider/pkg/influxdb/client/v2"
)

// HTTPClient HTTP Client
type HTTPClient struct {
	Client            client.Client
	BatchPointsConfig client.BatchPointsConfig
}

// FluxDBHttpWrite ...
func (p *HTTPClient) FluxDBHttpWrite(bp client.BatchPoints) (err error) {
	return p.Client.Write(bp)
}

// FluxDBHttpClose ...
func (p *HTTPClient) FluxDBHttpClose() (err error) {
	return p.Client.Close()
}
