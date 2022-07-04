package singleca

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/ztalab/cfssl/helpers"
)

func getTrustCerts(path string) ([]*x509.Certificate, error) {
	pemCerts, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("trust certificate file error: %v", err)
	}
	certs, err := helpers.ParseCertificatesPEM(pemCerts)
	if err != nil {
		return nil, fmt.Errorf("failed to get trust certificate: %v", err)
	}
	logger.Named("trust-certs").Infof("number of trust certificates obtained: %v", len(certs))
	return certs, nil
}
