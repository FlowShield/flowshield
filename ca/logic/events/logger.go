package events

import (
	"fmt"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	jsoniter "github.com/json-iterator/go"
)

const (
	LoggerName                = "events"
	CategoryWorkloadLifecycle = "workload_lifecycle"
)

var CategoriesStrings = map[string]string{
	CategoryWorkloadLifecycle: "Workload life cycle",
}

const (
	OperatorMSP = "MSP platform"
	OperatorSDK = "SDK"
)

type CertOp struct {
	UniqueId string `json:"unique_id"`
	SN       string `json:"sn"`
	AKI      string `json:"aki"`
}

// Op Operation record
type Op struct {
	Operator string      `json:"operator"` // Operator
	Category string      `json:"category"` // Classification
	Type     string      `json:"type"`     // Operation type
	Obj      interface{} `json:"obj"`      // Operation object
}

func (o *Op) Log() {
	objStr, _ := jsoniter.MarshalToString(o.Obj)
	logger.Named(LoggerName).
		With("flag", fmt.Sprintf("%s.%s", o.Category, o.Type)).
		With("data", o.Obj).
		Infof("Classification: %s, Operation: %s, Operator: %s, Operation object: %v", CategoriesStrings[o.Category], o.Type, o.Operator, objStr)
}
