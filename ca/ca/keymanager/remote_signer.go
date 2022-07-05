package keymanager

import (
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	jsoniter "github.com/json-iterator/go"
	cfssl_client "github.com/ztalab/cfssl/api/client"
	"github.com/ztalab/cfssl/cli/genkey"
	"github.com/ztalab/cfssl/csr"
	"github.com/ztalab/cfssl/signer"

	"github.com/cloudslit/cloudslit/ca/core"
)

// RemoteSigner ...
type RemoteSigner struct {
	logger *logger.Logger
}

// NewRemoteSigner ...
func NewRemoteSigner() *RemoteSigner {
	return &RemoteSigner{
		logger: logger.Named("remote-signer"),
	}
}

//Run calls the remote CA to sign the certificate and persist it
func (ss *RemoteSigner) Run() error {
	if core.Is.Config.Keymanager.SelfSign {
		return nil
	}
	key, cert, _ := GetKeeper().GetCachedSelfKeyPairPEM()
	if key != nil && cert != nil {
		ss.logger.Info("The certificate already exists. Skip the remote signing process")
		return nil
	}
	ss.logger.Warn("There is no certificate. You will sign the certificate remotely")
	g := &csr.Generator{Validator: genkey.Validator}
	csrBytes, key, err := g.ProcessRequest(getIntermediateCSRTemplate())
	if err != nil {
		ss.logger.Errorf("key, csr Production error: %v", err)
		return err
	}

	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Profile: "intermediate",
	}

	signReqBytes, _ := jsoniter.Marshal(&signReq)
	err = GetKeeper().RootClient.DoWithRetry(func(remote *cfssl_client.AuthRemote) error {
		certResp, err := remote.Sign(signReqBytes)
		if err != nil {
			return err
		}
		cert = certResp
		return nil
	})
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
