package keyprovider

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/flowshield/cfssl/csr"
	"github.com/flowshield/cfssl/helpers"
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/pkg/errors"

	"github.com/flowshield/flowshield/ca/pkg/keygen"
	"github.com/flowshield/flowshield/ca/pkg/spiffe"
)

const (
	DefaultPrivateKeyPath  = "/var/run/xkey/key.pem"
	DefaultCertificatePath = "/var/run/xkey/cert.pem"
)

// StandardPaths contains a path to a key file and certificate file.
type StandardPaths struct {
	KeyFile  string `json:"private_key"`
	CertFile string `json:"certificate"`
}

type xInternal struct {
	priv crypto.Signer
	cert *x509.Certificate

	mu sync.RWMutex

	// The PEM-encoded private key and certificate. This
	// is stored alongside the crypto.Signer and
	// x509.Certificate for convenience in marshaling and
	// calling tls.X509KeyPair directly.
	keyPEM  []byte
	certPEM []byte
}

// XKeyProvider provides unencrypted PEM-encoded certificates and
// private keys. If paths are provided, the key and certificate will
// be stored on disk.
type XKeyProvider struct {
	Paths               StandardPaths `json:"paths"`
	internal            xInternal
	*spiffe.IDGIdentity `json:"idg_identity"`
	DiskStore           bool
	logger              *logger.Logger
	CSRConf             keygen.CSRConf
}

// NewXKeyProvider sets up new XKeyProvider from the
// information contained in an Identity.
func NewXKeyProvider(id *spiffe.IDGIdentity) (*XKeyProvider, error) {
	if id == nil {
		return nil, errors.New("transport: the identity hasn't been initialised. Has it been loaded from disk?")
	}

	sp := &XKeyProvider{
		Paths: StandardPaths{
			KeyFile:  DefaultPrivateKeyPath,
			CertFile: DefaultCertificatePath,
		},
		internal: xInternal{
			mu: sync.RWMutex{},
		},
		IDGIdentity: id,
		logger:      logger.Named("keyprovider"),
	}

	err := sp.Check()
	if err != nil {
		return nil, err
	}

	return sp, nil
}

func (sp *XKeyProvider) resetCert() {
	sp.internal.mu.Lock()
	defer sp.internal.mu.Unlock()

	sp.internal.cert = nil
	sp.internal.certPEM = nil
}

func (sp *XKeyProvider) resetKey() {
	sp.internal.mu.Lock()
	defer sp.internal.mu.Unlock()

	sp.internal.priv = nil
	sp.internal.keyPEM = nil
}

var (
	// ErrMissingKeyPath is returned if the XKeyProvider has
	// specified a certificate path but not a key path.
	ErrMissingKeyPath = errors.New("transport: standard provider is missing a private key path to accompany the certificate path")

	// ErrMissingCertPath is returned if the XKeyProvider has
	// specified a private key path but not a certificate path.
	ErrMissingCertPath = errors.New("transport: standard provider is missing a certificate path to accompany the certificate path")
)

// Check ensures that the paths are valid for the provider.
func (sp *XKeyProvider) Check() error {
	if sp.Paths.KeyFile == "" && sp.Paths.CertFile == "" {
		return nil
	}

	if sp.Paths.KeyFile == "" {
		return ErrMissingKeyPath
	}

	if sp.Paths.CertFile == "" {
		return ErrMissingCertPath
	}

	return nil
}

// Persistent returns true if the key and certificate will be stored
// on disk.
func (sp *XKeyProvider) Persistent() bool {
	if sp.DiskStore && sp.Paths.KeyFile != "" && sp.Paths.CertFile != "" {
		return true
	}
	return false
}

// Generate generates a new private key.
func (sp *XKeyProvider) Generate(algo string, size int) (err error) {
	sp.resetKey()
	sp.resetCert()

	algo = strings.ToUpper(algo)
	priv, key, err := keygen.GenKey(keygen.SupportedSignatureAlgorithms(algo))
	if err != nil {
		return err
	}

	sp.internal.mu.Lock()
	defer sp.internal.mu.Unlock()

	sp.internal.priv = priv.(crypto.Signer)
	sp.internal.keyPEM = key

	sp.logger.Debugf("Create Private Key: %v", algo)

	return nil
}

// Certificate returns the associated certificate, or nil if
// one isn't ready.
func (sp *XKeyProvider) Certificate() *x509.Certificate {
	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()
	return sp.internal.cert
}

// CertificateRequest takes some metadata about a certificate request,
// and attempts to produce a certificate signing request suitable for
// sending to a certificate authority.
func (sp *XKeyProvider) CertificateRequest(_ *csr.CertificateRequest) ([]byte, error) {
	sp.internal.mu.RLock()
	if sp.internal.priv == nil {
		sp.internal.mu.RUnlock()
		err := sp.Generate(string(keygen.EcdsaSigAlg), 0)
		if err != nil {
			return nil, err
		}
	} else {
		sp.internal.mu.RUnlock()
	}
	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()

	return keygen.GenExtendWorkloadCSR(sp.internal.keyPEM, sp.IDGIdentity, sp.CSRConf)
}

// ErrCertificateUnavailable is returned when a key is available, but
// there is no accompanying certificate.
var ErrCertificateUnavailable = errors.New("transport: certificate unavailable")

// SetPrivateKeyPEM ...
func (sp *XKeyProvider) SetPrivateKeyPEM(pem []byte) error {
	sp.internal.mu.Lock()
	defer sp.internal.mu.Unlock()

	key, err := helpers.ParsePrivateKeyPEM(pem)
	if err != nil {
		return err
	}
	sp.internal.keyPEM = pem
	sp.internal.priv = key
	return nil
}

// Load a private key and certificate from disk.
func (sp *XKeyProvider) Load() (err error) {
	if !sp.Persistent() {
		return
	}

	var clearKey = true
	defer func() {
		if err != nil {
			if clearKey {
				sp.resetKey()
			}
			sp.resetCert()
		}
	}()

	sp.internal.keyPEM, err = ioutil.ReadFile(sp.Paths.KeyFile)
	if err != nil {
		return
	}

	sp.internal.priv, err = helpers.ParsePrivateKeyPEM(sp.internal.keyPEM)
	if err != nil {
		return
	}

	clearKey = false

	sp.internal.certPEM, err = ioutil.ReadFile(sp.Paths.CertFile)
	if err != nil {
		return ErrCertificateUnavailable
	}

	sp.internal.cert, err = helpers.ParseCertificatePEM(sp.internal.certPEM)
	if err != nil {
		err = errors.New("transport: invalid certificate")
		return
	}

	p, _ := pem.Decode(sp.internal.keyPEM)

	switch sp.internal.cert.PublicKey.(type) {
	case *rsa.PublicKey:
		if p.Type != "RSA PRIVATE KEY" {
			err = errors.New("transport: PEM type " + p.Type + " is invalid for an RSA key")
			return
		}
	case *ecdsa.PublicKey:
		if p.Type != "EC PRIVATE KEY" {
			err = errors.New("transport: PEM type " + p.Type + " is invalid for an ECDSA key")
			return
		}
	default:
		err = errors.New("transport: invalid public key type")
	}

	if err != nil {
		clearKey = true
		return
	}

	return nil
}

// Ready returns true if the provider has a key and certificate
// loaded. The certificate should be checked by the end user for
// validity.
func (sp *XKeyProvider) Ready() bool {
	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()

	switch {
	case sp.internal.priv == nil:
		return false
	case sp.internal.cert == nil:
		return false
	case sp.internal.keyPEM == nil:
		return false
	case sp.internal.certPEM == nil:
		return false
	default:
		return true
	}
}

// SetCertificatePEM receives a PEM-encoded certificate and loads it
// into the provider.
func (sp *XKeyProvider) SetCertificatePEM(certPEM []byte) error {
	sp.internal.mu.Lock()
	defer sp.internal.mu.Unlock()

	cert, err := helpers.ParseCertificatePEM(certPEM)
	if err != nil {
		return errors.New("transport: invalid certificate")
	}

	sp.internal.certPEM = certPEM
	sp.internal.cert = cert
	return nil
}

// SignalFailure is provided to implement the KeyProvider interface,
// and always returns false.
func (sp *XKeyProvider) SignalFailure(err error) bool {
	return false
}

// SignCSR takes a template certificate request and signs it.
func (sp *XKeyProvider) SignCSR(tpl *x509.CertificateRequest) ([]byte, error) {
	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()
	return x509.CreateCertificateRequest(rand.Reader, tpl, sp.internal.priv)
}

// Store writes the key and certificate to disk, if necessary.
func (sp *XKeyProvider) Store() error {
	if !sp.Ready() {
		return errors.New("transport: provider does not have a key and certificate")
	}

	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()

	return ioutil.WriteFile(sp.Paths.KeyFile, sp.internal.keyPEM, 0600)
}

// X509KeyPair returns a tls.Certificate for the provider.
func (sp *XKeyProvider) X509KeyPair() (tls.Certificate, error) {
	sp.internal.mu.RLock()
	defer sp.internal.mu.RUnlock()

	cert, err := tls.X509KeyPair(sp.internal.certPEM, sp.internal.keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	if cert.Leaf == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return tls.Certificate{}, err
		}
	}
	return cert, nil
}
