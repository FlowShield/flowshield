package initer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"github.com/cloudslit/cloudslit/provider/internal/config"
	"github.com/cloudslit/cloudslit/provider/pkg/certificate"
	"github.com/cloudslit/cloudslit/provider/pkg/util"
	"math/big"
	"net"
	"time"
)

const (
	TypeClient = "client"
	TypeServer = "server"
	TypeRelay  = "relay"
)

var (
	ErrCertParse = errors.New("certificate resolution error！")
	ErrCertType  = errors.New("sentinel type error！")
)

// certificate base field
type BasicCertConf struct {
	SiteID    string
	ClusterID string
	Type      string
}

func InitCert(certData []byte) (*BasicCertConf, map[string]interface{}, error) {
	p, _ := pem.Decode(certData)
	if p == nil {
		return nil, nil, ErrCertParse
	}
	cert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	basicConf := &BasicCertConf{}

	// parse attr
	mgr := certificate.New()
	attr, err := mgr.GetAttributesFromCert(cert)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	if t, ok := attr.Attrs["type"]; ok {
		if t, ok := t.(string); ok {
			basicConf.Type = t
		}
	}
	if !util.InArray(basicConf.Type, []string{TypeClient, TypeRelay, TypeServer}) {
		return nil, nil, ErrCertType
	}
	return basicConf, attr.Attrs, nil
}

// 自签证书
func InitSelfCert() error {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return err
	}
	subject := pkix.Name{
		Country:            []string{"CN"},
		Province:           []string{"BeiJing"},
		Organization:       []string{"Devops"},
		OrganizationalUnit: []string{"certDevops"},
		CommonName:         "127.0.0.1",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		return err
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})

	config.C.Certificate.CertPem = string(certPem)
	config.C.Certificate.KeyPem = string(keyPem)
	config.C.Certificate.CaPem = string(certPem)

	return nil
}
