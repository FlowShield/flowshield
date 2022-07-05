package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

type verifyCert struct {
	leaf          string
	intermediates []string
	roots         []string
	dnsName       string
}

// NewVerify Create a certificate validator
func NewVerify(cert, rootCert, dnsName string) *verifyCert {
	return &verifyCert{
		leaf:    cert,
		roots:   []string{rootCert},
		dnsName: dnsName,
	}
}

// expectAuthorityUnknown error handling
func (a *verifyCert) expectAuthorityUnknown(err error) error {
	e, ok := err.(x509.UnknownAuthorityError)
	if !ok {
		return errors.New("error was not UnknownAuthorityError: " + err.Error())
	}
	if e.Cert == nil {
		return errors.New("error was UnknownAuthorityError, but missing Cert: " + err.Error())
	}
	return err
}

// certificateFromPEM analytical certificate
func (a *verifyCert) certificateFromPEM(pemBytes string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemBytes))
	if block == nil {
		return nil, errors.New("failed to decode PEM")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (a *verifyCert) Verify() error {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		DNSName:       a.dnsName,
	}
	opts.Roots = x509.NewCertPool()
	for j, root := range a.roots {
		ok := opts.Roots.AppendCertsFromPEM([]byte(root))
		if !ok {
			return errors.New("failed to parse root #" + string(rune(j)))
		}
	}

	for j, intermediate := range a.intermediates {
		ok := opts.Intermediates.AppendCertsFromPEM([]byte(intermediate))
		if !ok {
			return errors.New("failed to parse intermediate #" + string(rune(j)))
		}
	}

	leaf, err := a.certificateFromPEM(a.leaf)
	if err != nil {
		return errors.New("failed to parse leaf:" + err.Error())
	}

	_, err = leaf.Verify(opts)
	if err != nil {
		return a.expectAuthorityUnknown(err)
	}
	return nil
}
