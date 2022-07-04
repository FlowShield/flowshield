package signature

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

// EcdsaSign The encryption result of the text signature is returned. The result is the serialization and splicing of the digital certificate R and s, and then converted into string with hex
func EcdsaSign(priv *ecdsa.PrivateKey, text []byte) (string, error) {
	hash := sha256.Sum256(text)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}
	return EcdsaSignEncode(r, s)
}

// EcdsaSignEncode r, s Convert to string
func EcdsaSignEncode(r, s *big.Int) (string, error) {
	rt, err := r.MarshalText()
	if err != nil {
		return "", err
	}
	st, err := s.MarshalText()
	if err != nil {
		return "", err
	}
	b := string(rt) + "," + string(st)
	return hex.EncodeToString([]byte(b)), nil
}

// EcdsaSignDecode r, s String parsing
func EcdsaSignDecode(sign string) (rint, sint big.Int, err error) {
	b, err := hex.DecodeString(sign)
	if err != nil {
		err = errors.New("decrypt error," + err.Error())
		return
	}
	rs := strings.Split(string(b), ",")
	if len(rs) != 2 {
		err = errors.New("decode fail")
		return
	}
	err = rint.UnmarshalText([]byte(rs[0]))
	if err != nil {
		err = errors.New("decrypt rint fail, " + err.Error())
		return
	}
	err = sint.UnmarshalText([]byte(rs[1]))
	if err != nil {
		err = errors.New("decrypt sint fail, " + err.Error())
		return
	}
	return
}

// EcdsaVerify Verify whether the text content is consistent with the signature. Use the public key to verify the signature and text content
func EcdsaVerify(text []byte, sign string, pubKey *ecdsa.PublicKey) (bool, error) {
	hash := sha256.Sum256(text)
	rint, sint, err := EcdsaSignDecode(sign)
	if err != nil {
		return false, err
	}
	result := ecdsa.Verify(pubKey, hash[:], &rint, &sint)
	return result, nil
}
