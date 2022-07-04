package keymanager

import (
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/ztalab/cfssl/initca"
)

// SelfSigner ...
type SelfSigner struct {
	logger *logger.Logger
}

// NewSelfSigner ...
func NewSelfSigner() *SelfSigner {
	return &SelfSigner{
		logger: logger.Named("self-signer"),
	}
}

// Run Self signed certificate and saved
func (ss *SelfSigner) Run() error {
	key, cert, _ := GetKeeper().GetCachedSelfKeyPairPEM()
	if key != nil && cert != nil {
		ss.logger.Info("The certificate already exists. Skip the self signing process")
		return nil
	}
	ss.logger.Warn("No certificate, self signed certificate")
	cert, _, key, err := initca.New(getRootCSRTemplate())
	if err != nil {
		ss.logger.Errorf("initca Create error: %v", err)
		return err
	}
	ss.logger.With("key", string(key), "cert", string(cert)).Debugf("Self signed certificate completed")
	if err = GetKeeper().SetKeyPairPEM(key, cert); err != nil {
		ss.logger.Errorf("Error saving certificate: %v", err)
		return err
	}

	return nil
}
