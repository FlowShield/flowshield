package signature

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"

	"github.com/pkg/errors"
)

// Signer ...
type Signer struct {
	priv crypto.PrivateKey
}

// NewSigner ...
func NewSigner(priv crypto.PrivateKey) *Signer {
	return &Signer{priv: priv}
}

// Sign
func (s *Signer) Sign(text []byte) (sign string, err error) {
	switch priv := s.priv.(type) {
	case *ecdsa.PrivateKey:
		sign, err = EcdsaSign(priv, text)
		return
	case *rsa.PrivateKey:
		// Todo supports RSA
		return "", errors.New("algo not supported")
	default:
		return "", errors.New("algo not supported")
	}
}

// Verifier ...
type Verifier struct {
	pub crypto.PublicKey
}

// NewVerifier ...
func NewVerifier(pub crypto.PublicKey) *Verifier {
	return &Verifier{pub: pub}
}

// Verify Verify signature
func (v *Verifier) Verify(text []byte, sign string) (bool, error) {
	switch pub := v.pub.(type) {
	case *ecdsa.PublicKey:
		return EcdsaVerify(text, sign, pub)
	case *rsa.PublicKey:
		// Todo supports RSA
	default:
	}
	return false, errors.New("algo not supported")
}
