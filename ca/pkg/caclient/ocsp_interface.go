package caclient

import "crypto/x509"

// OcspClient Ocsp Client
type OcspClient interface {
	Validate(leaf, issuer *x509.Certificate) (bool, error)
	Reset()
}
