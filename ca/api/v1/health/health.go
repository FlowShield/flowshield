package health

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/cloudSlit/cloudslit/ca/api/helper"
	"github.com/cloudSlit/cloudslit/ca/ca/keymanager"
	"github.com/cloudSlit/cloudslit/ca/core"
	cfClient "github.com/ztalab/cfssl/api/client"
)

// CfsslHealthAPI ...
const CfsslHealthAPI = "/api/v1/cfssl/health"

// HealthModule ...
type HealthModule struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Desc        string `json:"desc"`
	Message     string `json:"message"`
	State       int    `json:"state"`
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		IdleConnTimeout: 5 * time.Second,
	},
	Timeout: 3 * time.Second,
}

// Health ...
func Health(c *helper.HTTPWrapContext) (interface{}, error) {
	var hm []*HealthModule
	{
		// MySQL
		module := &HealthModule{
			Name:        "MySQL",
			DisplayName: "MySQL",
			State:       200,
		}
		if db, _ := core.Is.Db.DB(); db != nil {
			if err := db.Ping(); err != nil {
				module.Message = err.Error()
				module.State = 500
			}
		}
		hm = append(hm, module)
	}
	{
		// RootCA
		module := &HealthModule{
			Name:        "RootCA",
			DisplayName: "RootCA",
			State:       200,
		}
		keymanager.GetKeeper().RootClient.DoWithRetry(func(remote *cfClient.AuthRemote) error {
			caURL := remote.Hosts()[0]

			resp, err := httpClient.Get(caURL + CfsslHealthAPI)
			if err != nil {
				module.State = 500
				module.Message = err.Error()
			} else if resp.StatusCode >= 400 {
				module.Message = "response error"
				module.State = 500
			}
			return nil
		})
		hm = append(hm, module)
	}
	return hm, nil
}
