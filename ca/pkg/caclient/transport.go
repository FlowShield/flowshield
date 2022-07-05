package caclient

import (
	"crypto/tls"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/cloudflare/backoff"
	"github.com/ztalab/cfssl/csr"
	"github.com/ztalab/cfssl/transport/ca"
	"github.com/ztalab/cfssl/transport/core"
	"github.com/ztalab/cfssl/transport/kp"
	"github.com/ztalab/cfssl/transport/roots"
)

// A Transport is capable of providing transport-layer security using
// TLS.
type Transport struct {
	CertRefreshDurationRate int

	// Provider contains a key management provider.
	Provider kp.KeyProvider

	// CA contains a mechanism for obtaining signed certificates.
	CA ca.CertificateAuthority

	// TrustStore contains the certificates trusted by this
	// transport.
	TrustStore *roots.TrustStore

	// ClientTrustStore contains the certificate authorities to
	// use in verifying client authentication certificates.
	ClientTrustStore *roots.TrustStore

	// Identity contains information about the entity that will be
	// used to construct certificates.
	Identity *core.Identity

	// Backoff is used to control the behaviour of a Transport
	// when it is attempting to automatically update a certificate
	// as part of AutoUpdate.
	Backoff *backoff.Backoff

	// RevokeSoftFail, if true, will cause a failure to check
	// revocation (such that the revocation status of a
	// certificate cannot be checked) to not be treated as an
	// error.
	RevokeSoftFail bool

	manualRevoke bool

	logger *zap.SugaredLogger
}

// TLSClientAuthClientConfig Client TLS configuration, changing certificate dynamically
func (tr *Transport) TLSClientAuthClientConfig(host string) (*tls.Config, error) {
	return &tls.Config{
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			cert, err := tr.GetCertificate()
			if err != nil {
				tr.logger.Errorf("Client certificate acquisition error: %v", err)
				return nil, err
			}
			return cert, nil
		},
		RootCAs:      tr.TrustStore.Pool(),
		ServerName:   host,
		CipherSuites: core.CipherSuites,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// TLSClientAuthServerConfig The server TLS configuration needs to be changed dynamically
func (tr *Transport) TLSClientAuthServerConfig() (*tls.Config, error) {
	return &tls.Config{
		// Get configuration dynamically
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			tlsConfig := &tls.Config{
				GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
					cert, err := tr.GetCertificate()
					if err != nil {
						tr.logger.Errorf("Server certificate acquisition error: %v", err)
						return nil, err
					}
					return cert, nil
				},
				RootCAs:   tr.TrustStore.Pool(),
				ClientCAs: tr.ClientTrustStore.Pool(),
			}
			return tlsConfig, nil
		},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		CipherSuites: core.CipherSuites,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// TLSServerConfig is a general server configuration that should be
// used for non-client authentication purposes, such as HTTPS.
func (tr *Transport) TLSServerConfig() (*tls.Config, error) {
	return &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			cert, err := tr.GetCertificate()
			if err != nil {
				tr.logger.Errorf("Server certificate acquisition error: %v", err)
				return nil, err
			}
			return cert, nil
		},
		RootCAs:      tr.TrustStore.Pool(),
		ClientCAs:    tr.ClientTrustStore.Pool(),
		CipherSuites: core.CipherSuites,
		MinVersion:   tls.VersionTLS12,
		ClientAuth:   tls.VerifyClientCertIfGiven,
	}, nil
}

// Lifespan Returns the remaining replacement time of a certificate. If it is less than or equal to 0, the certificate must be replaced
// remain Total remaining time of certificate, ava update time
func (tr *Transport) Lifespan() (remain time.Duration, ava time.Duration) {
	cert := tr.Provider.Certificate()
	if cert == nil {
		return 0, 0
	}

	now := time.Now()
	if now.After(cert.NotAfter) {
		return 0, 0
	}
	remain = cert.NotAfter.Sub(now)

	certLong := cert.NotAfter.Sub(cert.NotBefore)
	ava = certLong / time.Duration(tr.CertRefreshDurationRate)

	if tr.manualRevoke {
		tr.manualRevoke = false
		return 0, 0
	}

	return remain, ava
}

// ManualRevoke ...
func (tr *Transport) ManualRevoke() {
	tr.manualRevoke = true
}

// RefreshKeys will make sure the Transport has loaded keys and has a
// valid certificate. It will handle any persistence, check that the
// certificate is valid (i.e. that its expiry date is within the
// Before date), and handle certificate reissuance as needed.
func (tr *Transport) RefreshKeys() (err error) {
	ch := make(chan error, 1)
	go func(tr *Transport) {
		ch <- tr.AsyncRefreshKeys()
	}(tr)
	select {
	case err := <-ch:
		return err
	case <-time.After(5 * time.Second): // 5 seconds timeout
		return errors.New("RefreshKeys timeout")
	}

}

// AsyncRefreshKeys timeout handler
func (tr *Transport) AsyncRefreshKeys() error {
	if !tr.Provider.Ready() {
		tr.logger.Debug("key and certificate aren't ready, loading")
		err := tr.Provider.Load()
		if err != nil && !errors.Is(err, kp.ErrCertificateUnavailable) {
			tr.logger.Debugf("failed to load keypair: %v", err)
			kr := tr.Identity.Request.KeyRequest
			if kr == nil {
				kr = csr.NewKeyRequest()
			}

			// Create a new private key
			tr.logger.Debug("Create a new private key")
			err = tr.Provider.Generate(kr.Algo(), kr.Size())
			if err != nil {
				tr.logger.Debugf("failed to generate key: %v", err)
				return err
			}
			tr.logger.Debug("Created successfully")
		}
	}

	// Certificate validity
	remain, lifespan := tr.Lifespan()
	if remain < lifespan || lifespan <= 0 {
		// Read the CSR configuration from the filled in request structure
		tr.logger.Debug("Create csr")
		req, err := tr.Provider.CertificateRequest(tr.Identity.Request)
		if err != nil {
			tr.logger.Debugf("couldn't get a CSR: %v", err)
			if tr.Provider.SignalFailure(err) {
				return tr.RefreshKeys()
			}
			return err
		}
		tr.logger.Debug("Create CSR complete")

		tr.logger.Debug("requesting certificate from CA")
		cert, err := tr.CA.SignCSR(req)
		if err != nil {
			if tr.Provider.SignalFailure(err) {
				return tr.RefreshKeys()
			}
			tr.logger.Debugf("failed to get the certificate signed: %v", err)
			return err
		}

		tr.logger.Debug("giving the certificate to the provider")
		err = tr.Provider.SetCertificatePEM(cert)
		if err != nil {
			tr.logger.Debugf("failed to set the provider's certificate: %v", err)
			if tr.Provider.SignalFailure(err) {
				return tr.RefreshKeys()
			}
			return err
		}

		if tr.Provider.Persistent() {
			tr.logger.Debug("storing the certificate")
			err = tr.Provider.Store()

			if err != nil {
				tr.logger.Debugf("the provider failed to store the certificate: %v", err)
				if tr.Provider.SignalFailure(err) {
					return tr.RefreshKeys()
				}
				return err
			}
		}
	}
	return nil
}

// GetCertificate ...
func (tr *Transport) GetCertificate() (*tls.Certificate, error) {
	tr.logger.Debug("keygen")
	if !tr.Provider.Ready() {
		tr.logger.Debug("transport isn't ready; attempting to refresh keypair")
		err := tr.RefreshKeys()
		if err != nil {
			tr.logger.Debugf("transport couldn't get a certificate: %v", err)
			return nil, err
		}
	}

	tr.logger.Debug("keypair")
	cert, err := tr.Provider.X509KeyPair()
	if err != nil {
		tr.logger.Debugf("couldn't generate an X.509 keypair: %v", err)
	}

	return &cert, nil
}

// AutoUpdate will automatically update the listener. If a non-nil
// certUpdates chan is provided, it will receive timestamps for
// reissued certificates. If errChan is non-nil, any errors that occur
// in the updater will be passed along.
func (tr *Transport) AutoUpdate() error {
	defer func() {
		if r := recover(); r != nil {
			tr.logger.Errorf("AutoUpdate certificates: %v", r)
		}
	}()
	remain, nextUpdateAt := tr.Lifespan()
	tr.logger.Debugf("attempting to refresh keypair")
	if remain > nextUpdateAt { // Failure to arrive at the rotation time: the rotation time is the certificate validity period of 1/2
		tr.logger.Debugf("Rotation time not reached %v %v", remain, nextUpdateAt)
		return nil
	}
	err := tr.RefreshKeys()
	if err != nil {
		retry := tr.Backoff.Duration()
		tr.logger.Debugf("failed to update certificate, will try again in %s", retry)
		return err
	}
	tr.logger.Debugf("certificate updated")
	tr.Backoff.Reset()
	return nil
}
