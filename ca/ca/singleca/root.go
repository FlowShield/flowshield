package singleca

import (
	"crypto"
	"crypto/x509"
	"strings"
	"time"

	"github.com/flowshield/cfssl/cli"
	// ...
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	// ...
	_ "github.com/flowshield/cfssl/cli/ocspsign"
	"github.com/flowshield/cfssl/ocsp"
	"github.com/flowshield/cfssl/signer"
	"github.com/flowshield/cfssl/signer/local"

	"github.com/flowshield/flowshield/ca/ca/keymanager"
	"github.com/flowshield/flowshield/ca/core"
)

var (
	conf       cli.Config
	s          signer.Signer
	ocspSigner ocsp.Signer
	db         *sqlx.DB
	router     = mux.NewRouter()
)

// registerHandlers instantiates various handlers and associate them to corresponding endpoints.
func registerHandlers() {
	logger := core.Is.Logger.Named("cfssl-handlers")

	disabled := make(map[string]bool)
	if conf.Disable != "" {
		for _, endpoint := range strings.Split(conf.Disable, ",") {
			disabled[endpoint] = true
		}
	}

	for path, getHandler := range endpoints {
		logger.Debugf("getHandler for %s", path)

		if _, ok := disabled[path]; ok {
			logger.Infof("endpoint '%s' is explicitly disabled", path)
		} else if handler, err := getHandler(); err != nil {
			logger.Warnf("endpoint '%s' is disabled: %v", path, err)
		} else {
			if path, handler, err = wrapHandler(path, handler, err); err != nil {
				logger.Warnf("endpoint '%s' is disabled by wrapper: %v", path, err)
			} else {
				logger.Infof("endpoint '%s' is enabled", path)
				router.Handle(path, handler)
			}
		}
	}
	logger.Info("Handler set up complete.")
}

func Server() (*mux.Router, error) {
	var err error
	logger := core.Is.Logger.Named("singleca")
	conf = cli.Config{
		Disable: "crl,gencrl,newcert,bundle,newkey,init_ca,scan,scaninfo,certinfo,ocspsign,/",
	}

	logger.Info("Initializing signer")

	// signer Assign to global variable s
	if s, err = local.NewDynamicSigner(
		func() crypto.Signer {
			priv, _, err := keymanager.GetKeeper().GetCachedSelfKeyPair()
			if err != nil {
				logger.Errorf("Error getting priv key: %v", err)
			}
			return priv
		}, func() *x509.Certificate {
			_, cert, err := keymanager.GetKeeper().GetCachedSelfKeyPair()
			if err != nil {
				logger.Errorf("Get cert error: %v", err)
			}
			return cert
		}, func() x509.SignatureAlgorithm {
			priv, _, err := keymanager.GetKeeper().GetCachedSelfKeyPair()
			if err != nil {
				logger.Errorf("Error getting priv key: %v", err)
			}
			return signer.DefaultSigAlgo(priv)
		}, core.Is.Config.Singleca.CfsslConfig.Signing); err != nil {
		logger.Errorf("couldn't initialize signer: %v", err)
		return nil, err
	}

	if ocspSigner, err = ocsp.NewDynamicSigner(
		func() *x509.Certificate {
			_, cert, err := keymanager.GetKeeper().GetCachedSelfKeyPair()
			if err != nil {
				logger.Errorf("Get cert error: %v", err)
			}
			return cert
		}, func() crypto.Signer {
			priv, _, err := keymanager.GetKeeper().GetCachedSelfKeyPair()
			if err != nil {
				logger.Errorf("Error getting priv key: %v", err)
			}
			return priv
		}, 4*24*time.Hour); err != nil {
		logger.Warnf("couldn't initialize ocsp signer: %v", err)
	}

	registerHandlers()

	return router, nil
}
