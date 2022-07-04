package caclient

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/pkg/errors"

	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/cloudSlit/cloudslit/ca/pkg/spiffe"
	"github.com/ztalab/cfssl/transport/core"
)

// TLSGenerator ...
type TLSGenerator struct {
	Cfg *tls.Config
}

// NewTLSGenerator ...
func NewTLSGenerator(cfg *tls.Config) *TLSGenerator {
	return &TLSGenerator{Cfg: cfg}
}

// ExtraValidator User defined verification function, which is executed after the certificate is verified successfully
type ExtraValidator func(identity *spiffe.IDGIdentity) error

// BindExtraValidator Register custom validation function
func (tg *TLSGenerator) BindExtraValidator(validator ExtraValidator) {
	vc := func(state tls.ConnectionState) error {
		// If there is no certificate, it will be blocked in the previous stage
		if len(state.PeerCertificates) == 0 {
			return nil
		}
		cert := state.PeerCertificates[0]
		var id *spiffe.IDGIdentity
		if len(cert.URIs) > 0 {
			id, _ = spiffe.ParseIDGIdentity(cert.URIs[0].String())
		}
		return validator(id)
	}
	getServerTls := tg.Cfg.GetConfigForClient
	if getServerTls != nil {
		// Server dynamic acquisition
		tg.Cfg.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			tlsCfg, err := getServerTls(info)
			if err != nil {
				return nil, err
			}
			tlsCfg.VerifyConnection = vc
			return tlsCfg, nil
		}
	} else {
		tg.Cfg.VerifyConnection = vc
	}
}

// TLSConfig Get golang native TLS config
func (tg *TLSGenerator) TLSConfig() *tls.Config {
	return tg.Cfg
}

// ClientTLSConfig ...
func (ex *Exchanger) ClientTLSConfig(host string) (*TLSGenerator, error) {
	lo := ex.logger
	lo.Debug("client tls started.")
	if _, err := ex.Transport.GetCertificate(); err != nil {
		return nil, errors.Wrap(err, "Client certificate acquisition error")
	}
	c, err := ex.Transport.TLSClientAuthClientConfig(host)
	if err != nil {
		return nil, err
	}
	c.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		if len(rawCerts) > 0 && len(verifiedChains) > 0 {
			leaf, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				lo.Errorf("leaf Certificate parsing error: %v", err)
				return err
			}
			if ok, err := ex.OcspFetcher.Validate(leaf, verifiedChains[0][1]); !ok {
				return err
			}
		}
		return nil
	}
	return NewTLSGenerator(c), nil
}

// ServerHTTPSConfig ...
func (ex *Exchanger) ServerHTTPSConfig() (*TLSGenerator, error) {
	lo := ex.logger
	lo.Debug("server tls started.")
	if _, err := ex.Transport.GetCertificate(); err != nil {
		return nil, errors.Wrap(err, "Server certificate acquisition error")
	}
	c, err := ex.Transport.TLSClientAuthServerConfig()
	if err != nil {
		return nil, err
	}
	c.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
		tlsConfig := &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := ex.Transport.GetCertificate()
				if err != nil {
					logger.Named("transport").Errorf("Server certificate acquisition error: %v", err)
					return nil, err
				}
				return cert, nil
			},
			ClientAuth:   tls.NoClientCert,
			CipherSuites: core.CipherSuites,
			MinVersion:   tls.VersionTLS12,
		}
		return tlsConfig, nil
	}
	return NewTLSGenerator(c), nil
}

// ServerTLSConfig ...
func (ex *Exchanger) ServerTLSConfig() (*TLSGenerator, error) {
	lo := ex.logger
	lo.Debug("server tls started.")
	if _, err := ex.Transport.GetCertificate(); err != nil {
		return nil, errors.Wrap(err, "Server certificate acquisition error")
	}
	c, err := ex.Transport.TLSClientAuthServerConfig()
	if err != nil {
		return nil, err
	}
	c.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
		tlsConfig := &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				cert, err := ex.Transport.GetCertificate()
				if err != nil {
					logger.Named("transport").Errorf("Server certificate acquisition error: %v", err)
					return nil, err
				}
				return cert, nil
			},
			RootCAs:   ex.Transport.TrustStore.Pool(),
			ClientCAs: ex.Transport.ClientTrustStore.Pool(),
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if len(rawCerts) > 0 && len(verifiedChains) > 0 {
					leaf, err := x509.ParseCertificate(rawCerts[0])
					if err != nil {
						lo.Errorf("leaf Certificate parsing error: %v", err)
						return err
					}
					if ok, err := ex.OcspFetcher.Validate(leaf, verifiedChains[0][1]); !ok {
						return err
					}
				}
				return nil
			},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			CipherSuites: core.CipherSuites,
			MinVersion:   tls.VersionTLS12,
		}
		return tlsConfig, nil
	}
	return NewTLSGenerator(c), nil
}
