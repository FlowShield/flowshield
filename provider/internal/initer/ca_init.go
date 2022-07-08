package initer

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/cloudslit/cloudslit/provider/pkg/certificate"
	"github.com/cloudslit/cloudslit/provider/pkg/util"
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
