package keymanager

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"math"
	"time"

	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudslit/cloudslit/ca/logic/schema"
	"github.com/cloudslit/cloudslit/ca/pkg/influxdb"
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/cloudslit/cloudslit/ca/pkg/memorycacher"
	"github.com/cloudslit/cloudslit/ca/pkg/vaultsecret"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	cfssl_client "github.com/ztalab/cfssl/api/client"
	"github.com/ztalab/cfssl/helpers"
	"github.com/ztalab/cfssl/hook"
	"github.com/ztalab/cfssl/info"
	"gorm.io/gorm"
)

// Keeper ...
type Keeper struct {
	DB         *gorm.DB
	cache      *memorycacher.Cache
	logger     *logger.Logger
	RootClient UpperClients
}

var (
	Std *Keeper
)

// ...
const (
	// SelfKeyPairName db row name
	SelfKeyPairName  = "ca"
	SelfKeyTrustName = "trust"
	// CacheKey
	cacheKeyPem  = "key-pem"
	cacheCertPem = "cert-pem"
	cacheKey     = "key"
	cacheCert    = "cert"
	// cacheTrustsPem = "trusts-pem"
	cacheTrusts = "trusts"
)

// InitKeeper ...
func InitKeeper() error {
	db := core.Is.Db
	var rootClients UpperClients
	var err error
	if !core.Is.Config.Keymanager.SelfSign {
		rootClients, err = NewUpperClients(core.Is.Config.Keymanager.UpperCa)
	}
	if err != nil {
		return errors.Wrap(err, "upper client Create error")
	}
	Std = &Keeper{
		DB:         db,
		logger:     logger.Named("keeper"),
		cache:      memorycacher.New(time.Hour, memorycacher.NoExpiration, math.MaxInt64),
		RootClient: rootClients,
	}
	return nil
}

// GetKeeper ...
func GetKeeper() *Keeper {
	defer func() {
		if err := recover(); err != nil {
			logger.Named("keeper").Fatal("Uninitialized")
		}
	}()
	return Std
}

// GetDBSelfKeyPairPEM ...
func (k *Keeper) GetDBSelfKeyPairPEM() (key, cert []byte, err error) {
	if !hook.EnableVaultStorage {
		keyPair := &model.SelfKeypair{}
		err = k.DB.Where("name = ?", SelfKeyPairName).Order("id desc").First(keyPair).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				k.logger.Warn("self Keys and certificates not found")
				return nil, nil, err
			}
			k.logger.Errorf("self-pair query error: %v", err)
			return nil, nil, err
		}
		if keyPair.PrivateKey.Valid {
			key = []byte(keyPair.PrivateKey.String)
		}
		if keyPair.Certificate.Valid {
			cert = []byte(keyPair.Certificate.String)
		}
	}

	if hook.EnableVaultStorage {
		certStr, keyStr, err := core.Is.VaultSecret.GetCertPEMKey(vaultsecret.CALocalStoreKey)
		if err != nil {
			k.logger.Errorf("vault Key and certificate read error: %s", err)
			return nil, nil, err
		}
		core.Is.Logger.With("key", keyStr, "cert", certStr).Debugf("Vault CA KEYPAIR")
		key = []byte(*keyStr)
		cert = []byte(*certStr)
	}

	return
}

// GetCachedTLSKeyPair ...
func (k *Keeper) GetCachedTLSKeyPair() (*tls.Certificate, error) {
	keyPEM, certPEM, err := k.GetCachedSelfKeyPairPEM()
	if err != nil {
		k.logger.Errorf("tls.Cert Get error： %v", err)
		return nil, err
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		k.logger.Errorf("tls.X509 error: %v", err)
		return nil, err
	}
	return &cert, nil
}

// GetCachedSelfKeyPair ...
func (k *Keeper) GetCachedSelfKeyPair() (key crypto.Signer, cert *x509.Certificate, err error) {
	if cachedKey, ok := k.cache.Get(cacheKey); ok {
		if v, ok := cachedKey.(crypto.Signer); ok {
			key = v
		}
	}
	if cachedCert, ok := k.cache.Get(cacheCert); ok {
		if v, ok := cachedCert.(*x509.Certificate); ok {
			cert = v
		}
	}

	if key != nil && cert != nil {
		return
	}

	keyPEM, certPEM, err := k.GetCachedSelfKeyPairPEM()
	if err != nil {
		k.logger.Errorf("Error getting cache keypair PEM: %v", err)
		return
	}
	priv, err := helpers.ParsePrivateKeyPEM(keyPEM)
	if err != nil {
		k.logger.With("key", string(keyPEM)).Errorf("Certificate key parsing error: %v", err)
		return
	}
	key = priv

	cert, err = helpers.ParseCertificatePEM(certPEM)
	if err != nil {
		k.logger.With("cert", string(certPEM)).Errorf("Certificate PEM parsing error: %v", err)
		return
	}

	k.cache.SetDefault(cacheKey, key)
	k.cache.SetDefault(cacheCert, cert)
	return
}

// GetCachedSelfKeyPairPEM ...
func (k *Keeper) GetCachedSelfKeyPairPEM() (key, cert []byte, err error) {
	if cachedKey, ok := k.cache.Get(cacheKeyPem); ok {
		if v, ok := cachedKey.([]byte); ok {
			key = v
		}
	}
	if cachedCert, ok := k.cache.Get(cacheCertPem); ok {
		if v, ok := cachedCert.([]byte); ok {
			cert = v
		}
	}
	if key != nil && cert != nil {
		return
	}
	key, cert, err = k.GetDBSelfKeyPairPEM()
	if key != nil && cert != nil {
		k.cache.SetDefault(cacheKeyPem, key)
		k.cache.SetDefault(cacheCertPem, cert)
	}
	return
}

// SetKeyPairPEM ...
func (k *Keeper) SetKeyPairPEM(key, cert []byte) error {
	keyPair := &model.SelfKeypair{
		Name:        SelfKeyPairName,
		PrivateKey:  sql.NullString{String: string(key), Valid: true},
		Certificate: sql.NullString{String: string(cert), Valid: true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if hook.EnableVaultStorage {
		keyPair.PrivateKey = sql.NullString{String: "", Valid: true}
		if err := core.Is.VaultSecret.StoreCertPEMKey(vaultsecret.CALocalStoreKey, string(cert), string(key)); err != nil {
			k.logger.Errorf("Vault write CA local store error: %s", err)
			return err
		}
	}
	if err := k.DB.Create(keyPair).Error; err != nil {
		k.logger.Errorf("Database insert error: %v", err)
		return err
	}
	k.cache.Flush()
	return nil
}

// GetL3CachedTrustCerts Memory > multi level cache > remote process > certificate
func (k *Keeper) GetL3CachedTrustCerts() (certs []*x509.Certificate, err error) {
	if cachedCerts, ok := k.cache.Get(cacheTrusts); ok {
		if v, ok := cachedCerts.([]*x509.Certificate); ok {
			return v, nil
		}
	}
	if !hook.EnableVaultStorage {
		dbTrustKeypair := &model.SelfKeypair{}
		dbErr := k.DB.Where("name = ?", SelfKeyTrustName).Order("id desc").First(dbTrustKeypair).Error
		if dbErr == nil {
			certs, err := helpers.ParseCertificatesPEM([]byte(dbTrustKeypair.Certificate.String))
			if err == nil {
				k.cache.SetDefault(cacheTrusts, certs)
				return certs, nil
			}
			k.logger.Errorf("DB Trust Certificate parsing error: %v", err)
		}

		if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
			k.logger.Errorf("DB get trust certificate error: %v", err)
		}
	}

	if hook.EnableVaultStorage {
		certsPEM, err := core.Is.VaultSecret.GetCertPEM(vaultsecret.CATructCertsKey)
		if err != nil {
			k.logger.Errorf("Vault get trust certificate error: %s", err)
		}
		certs, err := helpers.ParseCertificatesPEM([]byte(*certsPEM))
		if err == nil {
			k.cache.SetDefault(cacheTrusts, certs)
			return certs, nil
		}
		k.logger.Errorf("Vault Trust Certificate parsing error: %v", err)
	}

	certs, err = k.GetRemoteTrustCerts()
	if err != nil {
		k.logger.Errorf("Error getting trust certificate remotely: %v", err)
		return nil, err
	}
	if len(certs) > 0 {
		go func() {
			if err := k.saveTrustCerts(certs); err != nil {
				k.logger.Errorf("certs Storage error: %s", err)
			}
		}()
	}
	return certs, nil
}

// GetRemoteTrustCerts Obtain remote trust certificate (including root certificate and intermediate CA certificate)
func (k *Keeper) GetRemoteTrustCerts() (certs []*x509.Certificate, err error) {
	if core.Is.Config.Keymanager.SelfSign {
		return
	}
	reqBytes, _ := jsoniter.Marshal(&info.Req{
		Profile: "intermediate",
	})

	var resp *info.Resp
	err = k.RootClient.DoWithRetry(func(remote *cfssl_client.AuthRemote) error {
		infoResp, err := remote.Info(reqBytes)
		if err != nil {
			return err
		}
		if core.Is.Config.Influxdb.Enabled {
			core.Is.Metrics.AddPoint(&influxdb.MetricsData{
				Measurement: schema.MetricsUpperCaInfo,
				Fields: map[string]interface{}{
					"trust_certs_num": len(infoResp.TrustCertificates) + 1,
				},
				Tags: map[string]string{
					"type": schema.MetricsUpperCaTypeInfo,
					"host": schema.GetHostFromUrl(remote.Hosts()[0]),
				},
			})
		}
		resp = infoResp
		return nil
	})
	if err != nil {
		k.logger.Errorf("Error getting root certificate: %s", err)
		return nil, err
	}

	certsMap := make(map[string]*x509.Certificate, len(resp.TrustCertificates)+1)

	resp.TrustCertificates = append(resp.TrustCertificates, resp.Certificate)
	for _, certStr := range resp.TrustCertificates {
		cert, err := helpers.ParseCertificatePEM([]byte(certStr))
		if err != nil {
			k.logger.Errorf("ROOT Certificate parsing error: %v", err)
			return nil, err
		}
		certsMap[cert.SerialNumber.String()] = cert
	}

	for _, cert := range certsMap {
		certs = append(certs, cert)
	}

	return
}

func (k *Keeper) saveTrustCerts(certs []*x509.Certificate) error {
	certsPEM := helpers.EncodeCertificatesPEM(certs)
	trustKeypair := &model.SelfKeypair{
		Name:        SelfKeyTrustName,
		Certificate: sql.NullString{String: string(certsPEM), Valid: true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if hook.EnableVaultStorage {
		trustKeypair.Certificate = sql.NullString{String: "", Valid: true}
		if err := core.Is.VaultSecret.StoreCertPEM(vaultsecret.CATructCertsKey, string(certsPEM)); err != nil {
			k.logger.Errorf("vault Error saving trust certs: %s", err)
			return err
		}
	}
	// Insert here instead of update to ensure that there are records every time
	if err := k.DB.Create(trustKeypair).Error; err != nil {
		k.logger.Errorf("Database insert error: %v", err)
		return err
	}
	k.logger.With("num", len(certs)).Infof("Trust Insert certificate into database")
	return nil
}
