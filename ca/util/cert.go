package util

import "crypto/x509"

func GetSanURI(cert *x509.Certificate) string {
	if len(cert.URIs) > 0 {
		return cert.URIs[0].String()
	}
	return ""
}
