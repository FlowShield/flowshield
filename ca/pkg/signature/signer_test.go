package signature

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/flowshield/flowshield/ca/pkg/keygen"
	"testing"
)

func TestEcdsaSign(t *testing.T) {
	priv, _, _ := keygen.GenKey(keygen.EcdsaSigAlg)
	s := NewSigner(priv)
	sign, err := s.Sign([]byte("Test"))
	if err != nil {
		panic(err)
	}
	fmt.Println(sign)
}

func TestEcdsaVerify(t *testing.T) {
	text := []byte("Test")
	priv, _, _ := keygen.GenKey(keygen.EcdsaSigAlg)
	s := NewSigner(priv)
	sign, err := s.Sign(text)
	if err != nil {
		panic(err)
	}
	key := priv.(*ecdsa.PrivateKey)
	v := NewVerifier(key.Public())
	result, err := v.Verify(text, sign)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}
