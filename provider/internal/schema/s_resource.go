package schema

import (
	"net"
	"strconv"
	"strings"
)

type Resource struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
	Host string `json:"host"`
	Port string `json:"port"`
}

type Resources []*Resource

// Verify that the access resource exists
func (a Resources) VerifyResources(target Target) bool {
	isIP := false
	pip := net.ParseIP(target.Host)
	if pip != nil {
		isIP = true
	}
	for _, item := range a {
		if item.Host == "*" {
			return true
		}
		// dns validation
		if item.Type == "dns" && !isIP {
			if item.Host == target.Host &&
				item.CheckPort(target.Port) {
				return true
			}
			if item.CheckPort(target.Port) && strings.Index(item.Host, "*.") == 0 &&
				(strings.Contains(target.Host, item.Host[1:]) || item.Host[2:] == target.Host) {
				return true
			}
		}
		if item.Type == "cidr" && isIP && item.CheckPort(target.Port) {
			_, subnet, err := net.ParseCIDR(item.Host)
			if err != nil {
				// Resource restriction non-CIDR, direct comparison
				if target.Host == item.Host {
					return true
				}
				continue // The IP address does not match and is skipped
			}
			if subnet.Contains(pip) {
				return true
			}
		}
	}
	return false
}

// CheckPort Checking the Target Port eg:8080;9090;3000-4000
func (a *Resource) CheckPort(targetPort int) bool {
	fPort := strings.Split(a.Port, ";")
	for _, item := range fPort {
		if item == strconv.Itoa(targetPort) {
			return true
		}
		jPort := strings.Split(a.Port, "-")
		if len(jPort) == 2 {
			minPortInt, _ := strconv.Atoi(jPort[0])
			maxPortInt, _ := strconv.Atoi(jPort[1])
			if targetPort >= minPortInt && targetPort <= maxPortInt {
				return true
			}
		} else {
			continue
		}
	}
	return false
}
