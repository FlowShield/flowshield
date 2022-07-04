package revoke

import (
	"crypto/x509"
	"sync/atomic"
	"time"

	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/logic/schema"
	"github.com/cloudslit/cloudslit/ca/pkg/influxdb"
)

var overallRevokeCounter uint64

func CountAll() {
	if !core.Is.Config.Influxdb.Enabled {
		return
	}
	go func() {
		for {
			<-time.After(5 * time.Second)
			core.Is.Metrics.AddPoint(&influxdb.MetricsData{
				Measurement: schema.MetricsOverall,
				Fields: map[string]interface{}{
					"revoke_count": atomic.LoadUint64(&overallRevokeCounter),
				},
				Tags: map[string]string{
					"ip": schema.GetLocalIpLabel(),
				},
			})
		}
	}()
}

func AddMetricsPoint(cert *x509.Certificate) {
	if !core.Is.Config.Influxdb.Enabled {
		return
	}
	if cert == nil {
		return
	}
	atomic.AddUint64(&overallRevokeCounter, 1)
	core.Is.Metrics.AddPoint(&influxdb.MetricsData{
		Measurement: schema.MetricsCaRevoke,
		Fields: map[string]interface{}{
			"certs_num": 1,
		},
		Tags: map[string]string{
			"unique_id": cert.Subject.CommonName,
			"ip":        schema.GetLocalIpLabel(),
		},
	})
}
