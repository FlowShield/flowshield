package keygen

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"os"
	"strings"
	"time"

	cf_csr "github.com/flowshield/cfssl/csr"
	"github.com/flowshield/cfssl/helpers"

	"github.com/flowshield/flowshield/ca/pkg/pkiutil"
	"github.com/flowshield/flowshield/ca/pkg/spiffe"
	"github.com/flowshield/flowshield/ca/util"
)

type SupportedSignatureAlgorithms string

type KeySize int

const (
	EcdsaSigAlg SupportedSignatureAlgorithms = "ECDSA"
	RsaSigAlg   SupportedSignatureAlgorithms = "RSA"

	RsaKeySize2048  KeySize = 2048
	EcdsaKeySize256 KeySize = 256
)

// CSRConf custom csr config
type CSRConf struct {
	SNIHostnames []string
	IPAddresses  []string
}

type CertOptions struct {
	CN string

	// Comma-separated hostnames and IPs to generate a certificate for.
	// This can also be set to the identity running the workload,
	// like kubernetes service account.
	Host string

	// The NotBefore field of the issued certificate.
	NotBefore time.Time

	// TTL of the certificate. NotAfter - NotBefore.
	TTL time.Duration

	// Signer certificate.
	SignerCert *x509.Certificate

	// Signer private key.
	SignerPriv crypto.PrivateKey

	// Signer private key (PEM encoded).
	SignerPrivPem []byte

	// Organization for this certificate.
	Org string

	// Whether this certificate is used as signing cert for CA.
	IsCA bool

	// The type of Elliptical Signature algorithm to use
	// when generating private keys. Currently only ECDSA is supported.
	// If empty, RSA is used, otherwise ECC is used.
	SigAlg SupportedSignatureAlgorithms
}

// Generate Private Key
func GenKey(sigAlg SupportedSignatureAlgorithms) (priv interface{}, key []byte, err error) {
	var block pem.Block
	switch sigAlg {
	case RsaSigAlg:
		priv, err = (&cf_csr.KeyRequest{
			A: strings.ToLower(string(RsaSigAlg)),
			S: int(RsaKeySize2048),
		}).Generate()
		if err != nil {
			return nil, nil, err
		}
		keyBytes := x509.MarshalPKCS1PrivateKey(priv.(*rsa.PrivateKey))
		block = pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		}
	case EcdsaSigAlg:
		priv, err = (&cf_csr.KeyRequest{
			A: strings.ToLower(string(EcdsaSigAlg)),
			S: int(EcdsaKeySize256),
		}).Generate()
		if err != nil {
			return nil, nil, err
		}
		keyBytes, _ := x509.MarshalECPrivateKey(priv.(*ecdsa.PrivateKey))
		block = pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		}
	default:
		return nil, nil, errors.New("not supported private key algo")
	}

	key = pem.EncodeToMemory(&block)
	return priv, key, nil
}

// Generate CSR through key
// Support custom CSR requests
func GenCSR(key []byte, options CertOptions) ([]byte, error) {
	template, _ := pkiutil.GenCSRTemplate(pkiutil.CertOptions{
		Host:          options.Host,
		NotBefore:     options.NotBefore,
		TTL:           options.TTL,
		SignerCert:    options.SignerCert,
		SignerPriv:    options.SignerPriv,
		SignerPrivPem: options.SignerPrivPem,
		Org:           options.Org,
		IsCA:          options.IsCA,
		IsDualUse:     false,
	})
	template.Subject.CommonName = options.CN
	priv, err := helpers.ParsePrivateKeyPEM(key)
	if err != nil {
		return nil, err
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, crypto.PrivateKey(priv))
	if err != nil {
		return nil, err
	}
	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	csr := pem.EncodeToMemory(&block)
	return csr, nil
}

func GenWorkloadCSR(key []byte, id *spiffe.IDGIdentity) ([]byte, error) {
	hostname, _ := os.Hostname()
	ips := util.GetLocalIPs()
	hosts := make([]string, 0, 2+len(ips))
	hosts = append(hosts, id.String(), hostname)
	hosts = append(hosts, ips...)
	return GenCSR(key, CertOptions{
		Host: strings.Join(hosts, ","),
		Org:  id.ClusterID,
		CN:   id.UniqueID,
	})
}

// GenExtendWorkloadCSR Support custom CSR parameters
func GenExtendWorkloadCSR(key []byte, id *spiffe.IDGIdentity, csrConf CSRConf) ([]byte, error) {
	hostnames := make([]string, 0)
	if len(csrConf.SNIHostnames) > 0 {
		hostnames = append(hostnames, csrConf.SNIHostnames...)
	} else {
		hostname, _ := os.Hostname()
		hostnames = append(hostnames, hostname)
	}
	ips := util.GetLocalIPs()
	if len(csrConf.IPAddresses) > 0 {
		ips = csrConf.IPAddresses
	}
	hosts := make([]string, 0, 2+len(ips))
	hosts = append(hosts, id.String())
	hosts = append(hosts, ips...)
	hosts = append(hosts, hostnames...)
	return GenCSR(key, CertOptions{
		Host: strings.Join(hosts, ","),
		Org:  id.ClusterID,
		CN:   id.UniqueID,
	})
}

// GenCustomExtendCSR Generate business custom CSR with extended fields
func GenCustomExtendCSR(pemKey []byte, id *spiffe.IDGIdentity, opts *CertOptions, exts []pkix.Extension) ([]byte, error) {
	if opts.Host == "" {
		hostname, _ := os.Hostname()
		ips := util.GetLocalIPs()
		hosts := make([]string, 0, 2+len(ips))
		hosts = append(hosts, hostname)
		hosts = append(hosts, ips...)
		opts.Host = strings.Join(hosts, ",")
	}
	opts.Host += "," + id.String()
	opts.IsCA = false
	template, err := pkiutil.GenCSRTemplate(pkiutil.CertOptions{
		Host:          opts.Host,
		NotBefore:     opts.NotBefore,
		TTL:           opts.TTL,
		SignerCert:    opts.SignerCert,
		SignerPriv:    opts.SignerPriv,
		SignerPrivPem: opts.SignerPrivPem,
		Org:           id.ClusterID,
		IsCA:          opts.IsCA,
		IsDualUse:     false,
	})
	if err != nil {
		return nil, err
	}
	template.Subject.CommonName = opts.CN
	if len(exts) > 0 {
		template.ExtraExtensions = append(template.ExtraExtensions, exts...)
	}
	priv, err := helpers.ParsePrivateKeyPEM(pemKey)
	if err != nil {
		return nil, err
	}
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, template, crypto.PrivateKey(priv))
	if err != nil {
		return nil, err
	}
	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	csr := pem.EncodeToMemory(&block)
	return csr, nil
}
