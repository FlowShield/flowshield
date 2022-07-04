package ocsp

import (
	"sync/atomic"
	"time"

	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/cloudSlit/cloudslit/ca/logic/schema"
	"github.com/cloudSlit/cloudslit/ca/pkg/influxdb"
)

var (
	overallOcspSuccessCounter uint64
	overallOcspFailedCounter  uint64
	overallOcspCachedCounter  uint64
)

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
					"ocsp_success_count": atomic.LoadUint64(&overallOcspSuccessCounter),
				},
				Tags: map[string]string{
					"ip": schema.GetLocalIpLabel(),
				},
			})
			core.Is.Metrics.AddPoint(&influxdb.MetricsData{
				Measurement: schema.MetricsOverall,
				Fields: map[string]interface{}{
					"ocsp_failed_count": atomic.LoadUint64(&overallOcspFailedCounter),
				},
				Tags: map[string]string{
					"ip": schema.GetLocalIpLabel(),
				},
			})
			core.Is.Metrics.AddPoint(&influxdb.MetricsData{
				Measurement: schema.MetricsOverall,
				Fields: map[string]interface{}{
					"ocsp_cached_count": atomic.LoadUint64(&overallOcspCachedCounter),
				},
				Tags: map[string]string{
					"ip": schema.GetLocalIpLabel(),
				},
			})
		}
	}()
}

func AddMetricsPoint(uniqueID string, hitCache bool, certStatus string) {
	if !core.Is.Config.Influxdb.Enabled {
		return
	}
	cacheStatus := "miss"
	if hitCache {
		cacheStatus = "hit"
		atomic.AddUint64(&overallOcspCachedCounter, 1)
	}

	var fieldType string

	if certStatus == CertStatusGood {
		atomic.AddUint64(&overallOcspSuccessCounter, 1)
		fieldType = "success"
	} else {
		atomic.AddUint64(&overallOcspFailedCounter, 1)
		fieldType = "failed"
	}

	core.Is.Metrics.AddPoint(&influxdb.MetricsData{
		Measurement: schema.MetricsOcspResponses,
		Fields: map[string]interface{}{
			"times": 1,
		},
		Tags: map[string]string{
			"unique_id": uniqueID,
			"cache":     cacheStatus,
			"status":    certStatus,
			"type":      fieldType,
			"ip":        schema.GetLocalIpLabel(),
		},
	})
}
