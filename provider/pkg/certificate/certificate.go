package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/cloudslit/cloudslit/provider/pkg/util"
	"github.com/cloudslit/cloudslit/provider/pkg/util/json"
)

const (
	TypeClient = "client"
	TypeServer = "provider"
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
	NotBefore time.Time
	NotAfter  time.Time
}

func LoadCert(certData []byte) (*BasicCertConf, []byte, error) {
	p, _ := pem.Decode(certData)
	if p == nil {
		return nil, nil, ErrCertParse
	}
	cert, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	basicConf := &BasicCertConf{}
	basicConf.NotBefore = cert.NotBefore
	basicConf.NotAfter = cert.NotAfter
	// parse attr
	mgr := New()
	attr, err := mgr.GetAttributesFromCert(cert)
	if err != nil {
		return nil, nil, ErrCertParse
	}
	if t, ok := attr.Attrs["type"]; ok {
		if t, ok := t.(string); ok {
			basicConf.Type = t
		}
	}
	if !util.InArray(basicConf.Type, []string{TypeClient, TypeServer}) {
		return nil, nil, ErrCertType
	}
	str := json.MarshalToString(attr.Attrs)

	return basicConf, []byte(str), nil
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

	fmt.Println(string(certPem))
	fmt.Println(string(keyPem))

	return nil
}
