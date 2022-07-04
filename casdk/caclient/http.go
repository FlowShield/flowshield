package caclient

import (
	"crypto/tls"
	"net/http"
	"time"
)

var httpClient = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		},
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
	},
	Timeout: 1 * time.Second,
}
