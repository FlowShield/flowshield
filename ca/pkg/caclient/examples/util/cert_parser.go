package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cloudslit/cloudslit/ca/pkg/caclient"
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/ztalab/cfssl/helpers"
)

func ExtractCertFromExchanger(ex *caclient.Exchanger) {
	logger := logger.Named("keypair-exporter")
	tlsCert, err := ex.Transport.GetCertificate()
	if err != nil {
		logger.Errorf("TLS Certificate acquisition failed: %v", err)
		return
	}
	cert := helpers.EncodeCertificatePEM(tlsCert.Leaf)
	keyBytes, err := x509.MarshalPKCS8PrivateKey(tlsCert.PrivateKey)
	if err != nil {
		logger.Errorf("TLS certificate private key acquisition failed: %v", err)
		return
	}

	key := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	trustCerts := ex.Transport.TrustStore.Certificates()
	caCerts := make([][]byte, 0, len(trustCerts))

	fmt.Println("--- CA Certificate Stared ---")
	for _, caCert := range trustCerts {
		caCertBytes := helpers.EncodeCertificatePEM(caCert)
		caCerts = append(caCerts, caCertBytes)
		fmt.Println("---\n", string(caCertBytes), "\n---")
	}
	fmt.Println("--- CA Certificate End ---")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println()

	fmt.Println("--- Private key Stared ---\n", string(key), "\n--- Private key End ---")
	fmt.Println("--- Certificate Stared ---\n", string(cert), "\n--- Certificate End ---")
}
