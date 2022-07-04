package schema

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/mayocream/pki/pkg/x509util"
	"github.com/pkg/errors"
	cfCertInfo "github.com/ztalab/cfssl/certinfo"
	"github.com/ztalab/cfssl/helpers"

	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudslit/cloudslit/ca/pkg/caclient"
	"github.com/cloudslit/cloudslit/ca/pkg/spiffe"
)

func GetFullCertByX509Cert(cert *x509.Certificate) *FullCert {
	certBytes := helpers.EncodeCertificatePEM(cert)
	var certStr string
	if ctCert, err := x509util.CertificateFromPEM(certBytes); err == nil {
		certStr = x509util.CertificateToString(ctCert)
	}
	return &FullCert{
		SampleCert: SampleCert{
			SN:  cert.SerialNumber.String(),
			AKI: hex.EncodeToString(cert.SubjectKeyId),
			CN:  cert.Subject.CommonName,
			// TODO Join certificate ID acquisition role
			NotBefore: cert.NotBefore,
			Expiry:    cert.NotAfter,
		},
		CertInfo: ParseCertificate(cert),
		CertStr:  certStr,
		RawCert:  cert,
	}
}

func GetFullCertByModelCert(row *model.Certificates) (*FullCert, error) {
	cert, err := helpers.ParseCertificatePEM([]byte(row.Pem))
	if err != nil {
		return nil, errors.Wrap(err, "cert parse error")
	}
	var certStr string
	if ctCert, err := x509util.CertificateFromPEM([]byte(row.Pem)); err == nil {
		certStr = x509util.CertificateToString(ctCert)
	}
	return &FullCert{
		SampleCert: SampleCert{
			SN:        row.SerialNumber,
			AKI:       row.AuthorityKeyIdentifier,
			CN:        row.CommonName.String,
			Role:      caclient.Role(row.CaLabel.String),
			UniqueId:  row.CommonName.String,
			Status:    row.Status,
			IssuedAt:  row.IssuedAt,
			NotBefore: row.NotBefore,
			Expiry:    row.Expiry,
			RevokedAt: row.RevokedAt,
		},
		CertInfo: ParseCertificate(cert),
		CertStr:  certStr,
		RawCert:  cert,
	}, nil
}

// ParseCertificate parses an x509 certificate.
func ParseCertificate(cert *x509.Certificate) *Certificate {
	c := &Certificate{
		RawPEM:             string(helpers.EncodeCertificatePEM(cert)),
		SignatureAlgorithm: helpers.SignatureString(cert.SignatureAlgorithm),
		NotBefore:          cert.NotBefore,
		NotAfter:           cert.NotAfter,
		Subject:            cfCertInfo.ParseName(cert.Subject),
		Issuer:             cfCertInfo.ParseName(cert.Issuer),
		SANs:               cert.DNSNames,
		AKI:                formatKeyID(cert.AuthorityKeyId),
		SKI:                formatKeyID(cert.SubjectKeyId),
		SerialNumber:       cert.SerialNumber.String(),
	}
	for _, ip := range cert.IPAddresses {
		c.SANs = append(c.SANs, ip.String())
	}
	return c
}

func formatKeyID(id []byte) string {
	var s string

	for i, c := range id {
		if i > 0 {
			s += ":"
		}
		s += fmt.Sprintf("%02X", c)
	}

	return s
}

func GetCaMetadataFromX509Cert(cert *x509.Certificate) CaMetadata {
	var o, ou string
	if len(cert.Subject.Organization) > 0 {
		o = cert.Subject.Organization[0]
	}
	if len(cert.Subject.OrganizationalUnit) > 0 {
		ou = cert.Subject.OrganizationalUnit[0]
	}
	var id *spiffe.IDGIdentity
	if len(cert.URIs) > 0 {
		id, _ = spiffe.ParseIDGIdentity(cert.URIs[0].String())
	}
	return CaMetadata{
		O:        o,
		OU:       ou,
		SpiffeId: id,
	}
}
