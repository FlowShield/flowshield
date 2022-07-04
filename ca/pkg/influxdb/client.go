package influxdb

import (
	_ "github.com/cloudSlit/cloudslit/ca/pkg/influxdb/influxdb-client" // this is important because of the bug in go mod
	client "github.com/cloudSlit/cloudslit/ca/pkg/influxdb/influxdb-client/v2"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
)

// UDPClient UDP Client
type UDPClient struct {
	Conf              *Config
	BatchPointsConfig client.BatchPointsConfig
	client            client.Client
}

func (p *UDPClient) newUDPV1Client() *UDPClient {
	udpClient, err := client.NewUDPClient(client.UDPConfig{
		Addr: p.Conf.UDPAddress,
	})
	if err != nil {
		logger.Errorf("InfluxDBUDPClient err: %v", err)
	}
	p.client = udpClient
	return p
}

// FluxDBUDPWrite ...
func (p *UDPClient) FluxDBUDPWrite(bp client.BatchPoints) (err error) {
	err = p.newUDPV1Client().client.Write(bp)
	return
}

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
