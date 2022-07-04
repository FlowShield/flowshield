package caclient

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/hex"
	"net/http"

	"github.com/pkg/errors"

	"github.com/cloudSlit/cloudslit/ca/pkg/signature"
	jsoniter "github.com/json-iterator/go"
)

var revokePath = "/api/v1/cfssl/revoke"

// This type is meant to be unmarshalled from JSON
type RevokeRequest struct {
	Serial  string `json:"serial"`
	AKI     string `json:"authority_key_id"`
	Reason  string `json:"reason"`
	Nonce   string `json:"nonce"`
	Sign    string `json:"sign"`
	AuthKey string `json:"auth_key"`
	Profile string `json:"profile"`
}

// RevokeItSelf Revoke one's own certificate
func (ex *Exchanger) RevokeItSelf() error {
	tlsCert, err := ex.Transport.GetCertificate()
	if err != nil {
		return err
	}
	cert := tlsCert.Leaf
	priv := tlsCert.PrivateKey

	if err := revokeCert(ex.caAddr, priv, cert); err != nil {
		return err
	}
	ex.logger.With("sn", cert.SerialNumber.String()).Info("Service offline revoking its own certificate")

	return nil
}

func (cai *CAInstance) RevokeCert(priv crypto.PublicKey, cert *x509.Certificate) error {
	return revokeCert(cai.CaAddr, priv, cert)
}

func revokeCert(caAddr string, priv crypto.PublicKey, cert *x509.Certificate) error {
	s := signature.NewSigner(priv)

	nonce := cert.SerialNumber.String()

	sign, err := s.Sign([]byte(nonce))
	if err != nil {
		return err
	}

	req := &RevokeRequest{
		Serial: cert.SerialNumber.String(),
		AKI:    hex.EncodeToString(cert.AuthorityKeyId),
		Reason: "",
		Nonce:  nonce,
		Sign:   sign,
	}

	reqBytes, _ := jsoniter.Marshal(req)

	buf := bytes.NewBuffer(reqBytes)

	resp, err := httpClient.Post(caAddr+revokePath, "application/json", buf)
	if err != nil {
		return errors.Wrap(err, "Request error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Request error")
	}

	return nil
}
