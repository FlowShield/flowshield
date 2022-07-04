package ocsp

import (
	"encoding/hex"
	"math"
	"net/http"
	"time"

	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/cloudSlit/cloudslit/ca/pkg/memorycacher"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/helpers"
	"github.com/ztalab/cfssl/hook"
	"github.com/ztalab/cfssl/ocsp"
	"go.uber.org/zap"
	stdocsp "golang.org/x/crypto/ocsp"
	"gorm.io/gorm"

	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudSlit/cloudslit/ca/logic/events"
)

const (
	CertStatusGood           = "good"
	CertStatusUnknown        = "unknown"
	CertStatusNotFound       = "notfound"
	CertStatusServerError    = "servererror"
	CertStatusCertParseError = "certparseerror"
	CertStatusOCSPSignError  = "ocspsignerror"
)

var CertStatusIntMap = map[string]int{
	CertStatusGood:           200,
	CertStatusUnknown:        599,
	CertStatusNotFound:       404,
	CertStatusServerError:    500,
	CertStatusCertParseError: 400,
	CertStatusOCSPSignError:  502,
}

// SharedSources
type SharedSources struct {
	DB         *gorm.DB
	Cache      *memorycacher.Cache
	Logger     *zap.SugaredLogger
	OcspSigner ocsp.Signer
}

// NewSharedSources ...
func NewSharedSources(signer ocsp.Signer) (*SharedSources, error) {
	if core.Is.Db == nil {
		return nil, errors.New("database instance not found")
	}
	cacheTime := time.Duration(core.Is.Config.Ocsp.CacheTime)
	return &SharedSources{
		DB:         core.Is.Db,
		Logger:     logger.Named("ocsp-ss").SugaredLogger,
		Cache:      memorycacher.New(cacheTime*time.Minute, memorycacher.NoExpiration, math.MaxInt64),
		OcspSigner: signer,
	}, nil
}

// Response
func (ss *SharedSources) Response(req *stdocsp.Request) ([]byte, http.Header, error) {
	if req == nil {
		return nil, nil, errors.New("called with nil request")
	}

	aki := hex.EncodeToString(req.IssuerKeyHash)
	sn := req.SerialNumber

	if sn == nil {
		return nil, nil, errors.New("request contains no serial")
	}
	strSN := sn.String()

	if cachedResp, ok := ss.Cache.Get(strSN + aki); ok {
		if resp, ok := cachedResp.([]byte); ok {
			ss.Logger.With("sn", strSN, "aki", aki).Debugf("ocspResp cache")
			AddMetricsPoint("", true, CertStatusUnknown)
			return resp, nil, nil
		}
		ss.Logger.With("sn", strSN, "aki", aki).Errorf("cache Value parsing error")
	}

	// Database query
	certRecord := &model.Certificates{}
	if err := ss.DB.Where("serial_number = ? AND authority_key_identifier = ?", strSN, aki).First(certRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ss.Logger.With("sn", strSN, "aki", aki).Warnw("Certificate does not exist")
			AddMetricsPoint("", false, CertStatusNotFound)
			return nil, nil, ocsp.ErrNotFound
		}
		ss.Logger.With("sn", strSN, "aki", aki).Errorf("Certificate acquisition error: %v", err)
		AddMetricsPoint("", false, CertStatusServerError)
		return nil, nil, errors.Wrap(err, "server error")
	}

	if hook.EnableVaultStorage {
		pem, err := core.Is.VaultSecret.GetCertPEM(strSN)
		if err != nil {
			ss.Logger.With("sn", strSN, "aki", aki).Warnf("Vault Get error: %v", err)
		} else {
			certRecord.Pem = *pem
		}
	}

	cert, err := helpers.ParseCertificatePEM([]byte(certRecord.Pem))
	if err != nil {
		ss.Logger.With("sn", strSN, "aki", aki).Errorf("Certificate PEM parsing error: %v", err)
		AddMetricsPoint("", false, CertStatusCertParseError)
		return nil, nil, errors.Wrap(err, "cert err")
	}

	signReq := &ocsp.SignRequest{
		Certificate: cert,
		Status:      certRecord.Status,
		Reason:      int(certRecord.Reason.Int64),
		RevokedAt:   certRecord.RevokedAt,
	}

	ocspResp, err := ss.OcspSigner.Sign(*signReq)
	if err != nil {
		ss.Logger.With("sn", strSN, "aki", aki).Errorf("OCSP Sign error: %v", err)
		AddMetricsPoint(cert.Subject.CommonName, false, CertStatusOCSPSignError)
		return nil, nil, errors.Wrap(err, "internal err")
	}

	events.NewWorkloadLifeCycle("oscp-sign", events.OperatorSDK, events.CertOp{
		UniqueId: cert.Subject.CommonName,
		SN:       strSN,
		AKI:      aki,
	}).Log()

	ss.Cache.SetDefault(strSN+aki, ocspResp)

	ss.Logger.With("sn", strSN, "aki", aki).Infof("OCSP Signature Complete")

	AddMetricsPoint(cert.Subject.CommonName, false, CertStatusGood)
	return ocspResp, nil, nil
}
