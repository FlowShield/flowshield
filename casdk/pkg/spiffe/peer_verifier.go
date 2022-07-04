package spiffe

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// PeerCertVerifier is an instance to verify the peer certificate in the SPIFFE way using the retrieved root certificates.
type PeerCertVerifier struct {
	generalCertPool *x509.CertPool
	certPools       map[string]*x509.CertPool
}

// NewPeerCertVerifier returns a new PeerCertVerifier.
func NewPeerCertVerifier() *PeerCertVerifier {
	return &PeerCertVerifier{
		generalCertPool: x509.NewCertPool(),
		certPools:       make(map[string]*x509.CertPool),
	}
}

// GetGeneralCertPool returns generalCertPool containing all root certs.
func (v *PeerCertVerifier) GetGeneralCertPool() *x509.CertPool {
	return v.generalCertPool
}

// AddMapping adds a new trust domain to certificates mapping to the certPools map.
func (v *PeerCertVerifier) AddMapping(trustDomain string, certs []*x509.Certificate) {
	if v.certPools[trustDomain] == nil {
		v.certPools[trustDomain] = x509.NewCertPool()
	}
	for _, cert := range certs {
		v.certPools[trustDomain].AddCert(cert)
		v.generalCertPool.AddCert(cert)
	}
}

// AddMappingFromPEM adds multiple RootCA's to the spiffe Trust bundle in the trustDomain namespace
func (v *PeerCertVerifier) AddMappingFromPEM(trustDomain string, rootCertBytes []byte) error {
	block, rest := pem.Decode(rootCertBytes)
	var blockBytes []byte

	// Loop while there are no block are found
	for block != nil {
		blockBytes = append(blockBytes, block.Bytes...)
		block, rest = pem.Decode(rest)
	}

	rootCAs, err := x509.ParseCertificates(blockBytes)
	if err != nil {
		return fmt.Errorf("parse certificate from rootPEM got error: %v", err)
	}

	v.AddMapping(trustDomain, rootCAs)
	return nil
}

// AddMappings merges a trust domain to certs map to the certPools map.
func (v *PeerCertVerifier) AddMappings(certMap map[string][]*x509.Certificate) {
	for trustDomain, certs := range certMap {
		v.AddMapping(trustDomain, certs)
	}
}

// VerifyPeerCert is an implementation of tls.Config.VerifyPeerCertificate.
// It verifies the peer certificate using the root certificates associated with its trust domain.
func (v *PeerCertVerifier) VerifyPeerCert(rawCerts [][]byte, _ [][]*x509.Certificate) error {
	if len(rawCerts) == 0 {
		// Peer doesn't present a certificate. Just skip. Other authn methods may be used.
		return nil
	}
	var peerCert *x509.Certificate
	intCertPool := x509.NewCertPool()
	for id, rawCert := range rawCerts {
		cert, err := x509.ParseCertificate(rawCert)
		if err != nil {
			return err
		}
		if id == 0 {
			peerCert = cert
		} else {
			intCertPool.AddCert(cert)
		}
	}
	if peerCert == nil {
		return fmt.Errorf("peer certificate not found")
	}
	if len(peerCert.URIs) != 1 {
		return fmt.Errorf("peer certificate does not contain 1 URI type SAN, detected %d", len(peerCert.URIs))
	}
	uri := peerCert.URIs[0].String()
	id, err := ParseIDGIdentity(uri)
	if err != nil {
		return err
	}
	trustDomain := id.SiteID
	rootCertPool, ok := v.certPools[trustDomain]
	if !ok {
		return fmt.Errorf("no cert pool found for trust domain %s", trustDomain)
	}

	_, err = peerCert.Verify(x509.VerifyOptions{
		Roots:         rootCertPool,
		Intermediates: intCertPool,
	})
	return err
}
