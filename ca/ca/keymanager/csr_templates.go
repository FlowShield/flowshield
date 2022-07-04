package keymanager

import (
	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/ztalab/cfssl/csr"
)

// getRootCSRTemplate Root CA
var getRootCSRTemplate = func() *csr.CertificateRequest {
	return &csr.CertificateRequest{
		Names: []csr.Name{
			{O: core.Is.Config.Keymanager.CsrTemplates.RootCa.O},
		},
		KeyRequest: &csr.KeyRequest{
			A: "rsa",
			S: 4096,
		},
		CA: &csr.CAConfig{
			Expiry: core.Is.Config.Keymanager.CsrTemplates.RootCa.Expiry,
		},
	}
}

// getIntermediateCSRTemplate
var getIntermediateCSRTemplate = func() *csr.CertificateRequest {
	return &csr.CertificateRequest{
		Names: []csr.Name{
			{
				O:  core.Is.Config.Keymanager.CsrTemplates.IntermediateCa.O,
				OU: core.Is.Config.Keymanager.CsrTemplates.IntermediateCa.Ou,
			},
		},
		KeyRequest: &csr.KeyRequest{
			A: "rsa",
			S: 4096,
		},
		CA: &csr.CAConfig{
			Expiry: core.Is.Config.Keymanager.CsrTemplates.IntermediateCa.Expiry,
		},
	}
}
