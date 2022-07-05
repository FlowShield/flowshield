package caclient

import (
	"github.com/cloudflare/backoff"
	"github.com/cloudslit/cloudslit/ca/pkg/keyprovider"
	"github.com/cloudslit/cloudslit/ca/pkg/spiffe"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/hook"
	"github.com/ztalab/cfssl/transport"
	"github.com/ztalab/cfssl/transport/roots"
	"go.uber.org/zap"
	"net/url"
	"reflect"
)

const (
	// CertRefreshDurationRate Certificate cycle time rate
	CertRefreshDurationRate int = 2
)

// Exchanger ...
type Exchanger struct {
	Transport   *Transport
	IDGIdentity *spiffe.IDGIdentity
	OcspFetcher OcspClient

	caAddr string
	logger *zap.SugaredLogger

	caiConf *Conf
}

func init() {
	// Cfssl API client connects to API server without certificate verification (one-way TLS)
	hook.ClientInsecureSkipVerify = true
}

// NewExchangerWithKeypair ...
func (cai *CAInstance) NewExchangerWithKeypair(id *spiffe.IDGIdentity, keyPEM []byte, certPEM []byte) (*Exchanger, error) {
	tr, err := cai.NewTransport(id, keyPEM, certPEM)
	if err != nil {
		return nil, err
	}
	of, err := NewOcspMemCache(cai.Logger.Sugar().Named("ocsp"), cai.Conf.OcspAddr)
	if err != nil {
		return nil, err
	}
	return &Exchanger{
		Transport:   tr,
		IDGIdentity: id,
		OcspFetcher: of,
		logger:      cai.Logger.Sugar().Named("ca"),
		caAddr:      cai.CaAddr,

		caiConf: &cai.Conf,
	}, nil
}

// NewExchanger ...
func (cai *CAInstance) NewExchanger(id *spiffe.IDGIdentity) (*Exchanger, error) {
	tr, err := cai.NewTransport(id, nil, nil)
	if err != nil {
		return nil, err
	}
	of, err := NewOcspMemCache(cai.Logger.Sugar().Named("ocsp"), cai.Conf.OcspAddr)
	if err != nil {
		return nil, err
	}
	return &Exchanger{
		Transport:   tr,
		IDGIdentity: id,
		OcspFetcher: of,
		logger:      cai.Logger.Sugar().Named("ca"),
		caAddr:      cai.CaAddr,

		caiConf: &cai.Conf,
	}, nil
}

// NewTransport ...
func (cai *CAInstance) NewTransport(id *spiffe.IDGIdentity, keyPEM []byte, certPEM []byte) (*Transport, error) {
	l := cai.Logger.Sugar()

	l.Debug("NewTransport Start")

	if _, err := url.Parse(cai.CaAddr); err != nil {
		return nil, errors.Wrap(err, "CA ADDR Error")
	}

	var tr = &Transport{
		CertRefreshDurationRate: CertRefreshDurationRate,
		Identity:                cai.CFIdentity,
		Backoff:                 &backoff.Backoff{},
		logger:                  l.Named("ca"),
	}

	l.Debugf("[NEW]: Certificate rotation rate: %v", tr.CertRefreshDurationRate)

	l.Debug("roots Initialization")
	store, err := roots.New(cai.CFIdentity.Roots)
	if err != nil {
		return nil, err
	}
	tr.TrustStore = store

	l.Debug("client roots Initialization")
	if len(cai.CFIdentity.ClientRoots) > 0 {
		if !reflect.DeepEqual(cai.CFIdentity.Roots, cai.CFIdentity.ClientRoots) {
			store, err = roots.New(cai.CFIdentity.ClientRoots)
			if err != nil {
				return nil, err
			}
		}

		tr.ClientTrustStore = store
	}

	l.Debug("xkeyProvider Initialization")
	xkey, err := keyprovider.NewXKeyProvider(id)
	if err != nil {
		return nil, err
	}

	xkey.CSRConf = cai.CSRConf
	if keyPEM != nil && certPEM != nil {
		l.Debug("xkeyProvider set up keyPEM")
		if err := xkey.SetPrivateKeyPEM(keyPEM); err != nil {
			return nil, err
		}
		l.Debug("xkeyProvider set up certPEM")
		if err := xkey.SetCertificatePEM(certPEM); err != nil {
			return nil, err
		}
	}
	tr.Provider = xkey

	l.Debug("CA Initialization")
	tr.CA, err = transport.NewCA(cai.CFIdentity)
	if err != nil {
		return nil, err
	}

	return tr, nil
}

// RotateController ...
func (ex *Exchanger) RotateController() *RotateController {
	return &RotateController{
		transport:   ex.Transport,
		rotateAfter: ex.caiConf.RotateAfter,
		logger:      ex.logger.Named("rotator"),
	}
}
