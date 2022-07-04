package schema

import (
	"net"
	"net/url"

	"github.com/cloudslit/cloudslit/ca/core/config"
)

const (
	MetricsOcspResponses = config.MetricsTablePrefix + "ocsp_responses"
	MetricsUpperCaInfo   = config.MetricsTablePrefix + "upper_ca_info"
	MetricsOverall       = config.MetricsTablePrefix + "ca_overall"
	MetricsCaSign        = config.MetricsTablePrefix + "ca_sign"
	MetricsCaRevoke      = config.MetricsTablePrefix + "ca_revoke"
	MetricsCaCpuMem      = config.MetricsTablePrefix + "cpu_mem"

	MetricsUpperCaTypeInfo = "ca_info"

	MetricsLabelIp = "ip"
)

func GetHostFromUrl(addr string) string {
	host, err := url.Parse(addr)
	if err != nil {
		return ""
	}
	return host.Host
}

func GetLocalIpLabel() string {
	return internetIP
}

var internetIP = getInternetIP()

func getInternetIP() (IP string) {
	// Find native IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				if ip4[0] == 10 {
					// Assign new IP
					IP = ip4.String()
				}
			}
		}
	}
	return
}
