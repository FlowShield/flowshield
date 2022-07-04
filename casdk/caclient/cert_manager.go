package caclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"net/http"

	"github.com/ztalab/cfssl/hook"

	"github.com/ztalab/cfssl/helpers"
	"github.com/ztalab/cfssl/info"

	jsoniter "encoding/json"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/api/client"
	"github.com/ztalab/cfssl/auth"
	"github.com/ztalab/cfssl/signer"
	"go.uber.org/zap"
)

// CertManager Certificate manager
type CertManager struct {
	logger      *zap.SugaredLogger
	apiClient   *client.AuthRemote
	profile     string
	caAddr      string
	authKey     string
	ocspFetcher OcspClient
	// TODO Certificate storage
	caCertTmp *x509.Certificate
}

// NewCertManager Create certificate management Instance
func (cai *CAInstance) NewCertManager() (*CertManager, error) {
	ap, err := auth.New(cai.Conf.CFIdentity.Profiles["cfssl"]["auth-key"], nil)
	if err != nil {
		return nil, errors.Wrap(err, "Auth key Configuration error")
	}
	caAddr := cai.CaAddr
	ocspAddr := cai.OcspAddr
	apiClient := client.NewAuthServer(caAddr, &tls.Config{
		InsecureSkipVerify: true, //nolint:gosec
	}, ap)
	profile := cai.Conf.CFIdentity.Profiles["cfssl"]["profile"]
	if profile == "" {
		return nil, errors.New("profile could not be empty")
	}
	ocspFetcher, err := NewOcspMemCache(cai.Logger.Sugar().Named("ocsp"), ocspAddr)
	if err != nil {
		return nil, err
	}
	cm := &CertManager{
		logger:      cai.Logger.Sugar().Named("cert-manager"),
		apiClient:   apiClient,
		profile:     profile,
		caAddr:      caAddr,
		ocspFetcher: ocspFetcher,
		authKey:     cai.Conf.CFIdentity.Profiles["cfssl"]["auth-key"],
	}

	cm.caCertTmp, err = cm.CACert()
	if err != nil {
		return nil, err
	}

	return cm, nil
}

// SignPEM ...
func (cm *CertManager) SignPEM(csrPEM []byte, uniqueID string) ([]byte, error) {
	if csrPEM == nil {
		return nil, errors.New("empty input")
	}

	signReq := signer.SignRequest{
		Request: string(csrPEM),
		Profile: cm.profile,
		Metadata: map[string]interface{}{
			hook.MetadataUniqueID: uniqueID,
		},
	}

	csr, err := helpers.ParseCSRPEM(csrPEM)
	if err != nil {
		return nil, err
	}
	signReq.Subject = &signer.Subject{
		CN: csr.Subject.CommonName,
	}

	signReqBytes, err := jsoniter.Marshal(&signReq)
	if err != nil {
		return nil, err
	}

	cm.logger.With("req", signReq).Debug("Request for certificate")

	certPEM, err := cm.apiClient.Sign(signReqBytes)
	if err != nil {
		cm.logger.Errorf("Request to issue certificate failed: %s", err)
		return nil, err
	}

	return certPEM, nil
}

// RevokeIDGRegistryCert ...
func (cm *CertManager) RevokeIDGRegistryCert(certPEM []byte) error {
	cert, err := helpers.ParseCertificatePEM(certPEM)
	if err != nil {
		return err
	}

	req := &RevokeRequest{
		Serial:  cert.SerialNumber.String(),
		AKI:     hex.EncodeToString(cert.AuthorityKeyId),
		Reason:  "", // Default to 0
		AuthKey: cm.authKey,
		Profile: cm.profile,
	}

	reqBytes, _ := jsoniter.Marshal(req)

	buf := bytes.NewBuffer(reqBytes)

	resp, err := httpClient.Post(cm.caAddr+revokePath, "application/json", buf)
	if err != nil {
		return errors.Wrap(err, "Request error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		cm.logger.With("status", resp.StatusCode).Errorf("Request error")
		return errors.New("Request error")
	}

	return nil
}

// RevokeByKeyPEM ...
func (cm *CertManager) RevokeByKeyPEM(keyPEM, certPEM []byte) error {
	if keyPEM == nil || certPEM == nil {
		return errors.New("empty input")
	}
	priv, err := helpers.ParsePrivateKeyPEM(keyPEM)
	if err != nil {
		return err
	}
	cert, err := helpers.ParseCertificatePEM(certPEM)
	if err != nil {
		return err
	}
	return revokeCert(cm.caAddr, priv, cert)
}

// VerifyCertDefaultIssuer ...
func (cm *CertManager) VerifyCertDefaultIssuer(leafPEM []byte) error {
	leaf, err := helpers.ParseCertificatePEM(leafPEM)
	if err != nil {
		return err
	}
	ok, err := cm.ocspFetcher.Validate(leaf, cm.caCertTmp)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Certificate revoked")
	}
	return nil
}

// CACert ...
func (cm *CertManager) CACert() (*x509.Certificate, error) {
	reqBytes, _ := jsoniter.Marshal(&info.Req{
		Profile: cm.profile,
	})
	resp, err := cm.apiClient.Info(reqBytes)
	if err != nil {
		return nil, err
	}
	cert, err := helpers.ParseCertificatePEM([]byte(resp.Certificate))
	if err != nil {
		return nil, err
	}
	return cert, nil
}

// CACertsPEM ...
func (cm *CertManager) CACertsPEM() ([]byte, error) {
	reqBytes, _ := jsoniter.Marshal(&info.Req{
		Profile: cm.profile,
	})
	resp, err := cm.apiClient.Info(reqBytes)
	if err != nil {
		return nil, err
	}
	caCert, err := helpers.ParseCertificatePEM([]byte(resp.Certificate))
	if err != nil {
		return nil, err
	}
	certs := make([]*x509.Certificate, 0)
	for _, trustCert := range resp.TrustCertificates {
		cert, err := helpers.ParseCertificatePEM([]byte(trustCert))
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}
	certs = append(certs, caCert)

	certsPem := helpers.EncodeCertificatesPEM(certs)

	return certsPem, nil
}
