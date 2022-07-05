package schema

import (
	"github.com/cloudslit/cloudslit/provider/pkg/errors"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
)

type RelayConfig struct {
	Type string `json:"type"`
	Port int    `json:"port"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func ParseRelayConfig(attrs map[string]interface{}) (*RelayConfig, error) {
	var result RelayConfig
	attrByte, err := json.Marshal(attrs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = json.Unmarshal(attrByte, &result)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &result, nil
}
