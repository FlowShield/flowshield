package metrics

import (
	"context"
	"github.com/cloudSlit/cloudslit/provider/internal/config"
	"github.com/cloudSlit/cloudslit/provider/pkg/errors"
	"github.com/cloudSlit/cloudslit/provider/pkg/influxdb"
	"github.com/cloudSlit/cloudslit/provider/pkg/logger"
)

const (
	ReqSuccess = "success"
	ReqFail    = "fail"

	Prefix     = "za-sentinel"
	MetricsReq = Prefix + "req"
)

type Metrics struct {
	PodIP       string `json:"pod_ip"`
	UniqueID    string `json:"unique_id"`
	Hostname    string `json:"hostname"`
	ServiceName string `json:"service_name"`
	Delay       string `json:"delay"`
	Status      string `json:"status"`
	Operator    string `json:"operator"`
}

func AddDelayPoint(ctx context.Context, operator, status, delay, id, name string) {
	if !config.C.Influxdb.Enabled {
		return
	}
	fields := make(map[string]interface{})
	fields["delay"] = delay
	fields["status"] = status

	tags := make(map[string]string)
	tags["pod_ip"] = config.C.Common.PodIP
	tags["unique_id"] = config.C.Common.UniqueID
	tags["hostname"] = config.C.Common.Hostname
	tags["app_name"] = config.C.Common.AppName
	tags["operator"] = operator
	tags["id"] = id
	tags["name"] = name

	err := config.Is.Metrics.AddPoint(&influxdb.MetricsData{
		Measurement: MetricsReq,
		Fields:      fields,
		Tags:        tags,
	})
	if err != nil {
		logger.WithErrorStack(ctx, errors.WithStack(err)).Errorf("Failed to add sequence logs. Procedure：%v", err)
	}
}
