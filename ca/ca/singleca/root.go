package singleca

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	// ...
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/certdb/sql"
	"github.com/ztalab/cfssl/cli"
	// ...
	_ "github.com/ztalab/cfssl/cli/ocspsign"
	"github.com/ztalab/cfssl/ocsp"
	"github.com/ztalab/cfssl/signer"
	"github.com/ztalab/cfssl/signer/local"

	"github.com/cloudSlit/cloudslit/ca/ca/keymanager"
	ocsp_responder "github.com/cloudSlit/cloudslit/ca/ca/ocsp"
	"github.com/cloudSlit/cloudslit/ca/ca/upperca"
	"github.com/cloudSlit/cloudslit/ca/core"
)

var (
	conf        cli.Config
	s           signer.Signer
	ocspSigner  ocsp.Signer
	db          *sqlx.DB
	router      = mux.NewRouter()
	proxyClient = resty.NewWithClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
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

	// Certificate signature
	if core.Is.Config.Keymanager.SelfSign {
		conf = cli.Config{
			Disable: "sign,crl,gencrl,newcert,bundle,newkey,init_ca,scan,scaninfo,certinfo,ocspsign,/",
		}
		if err := keymanager.NewSelfSigner().Run(); err != nil {
			logger.Fatalf("Self signed certificate error: %v", err)
		}
		router.PathPrefix("/api/v1/cap/").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			localPort := core.Is.Config.HTTP.Listen
			splits := strings.Split(localPort, ":")
			localPort = splits[len(splits)-1]
			localUrl := "http://127.0.0.1:" + localPort + "/api/v1/"
			localUrl += strings.TrimPrefix(strings.Replace(request.URL.Path, "/api/v1/cap/", "", 1), "/")
			var resp *resty.Response
			var err error
			switch request.Method {
			case http.MethodGet:
				resp, err = proxyClient.R().
					Get(localUrl)
			case http.MethodPost:
				bodyBytes, _ := ioutil.ReadAll(request.Body)
				resp, err = proxyClient.R().
					SetBody(bodyBytes).
					Post(localUrl)
			default:
				writer.WriteHeader(404)
				writer.Write([]byte("404 not found"))
			}

			if err != nil {
				logger.Errorf("Request error: %s", err)
				writer.WriteHeader(500)
				writer.Write([]byte("server error"))
			}

			writer.WriteHeader(200)
			writer.Write(resp.Body())
		})
	} else {
		conf = cli.Config{
			Disable: "crl,gencrl,newcert,bundle,newkey,init_ca,scan,scaninfo,certinfo,ocspsign,/",
		}
		if err := keymanager.NewRemoteSigner().Run(); err != nil {
			logger.Fatalf("Remote signing certificate error: %v", err)
		}
		// Superior CA health check
		go upperca.NewChecker().Run()
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
	db, err = sqlx.Open("mysql", core.Is.Config.Mysql.Dsn)
	if err != nil {
		logger.Errorf("Sqlx Initialization error: %v", err)
		return nil, err
	}
	s.SetDBAccessor(sql.NewAccessor(db))

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

	endpoints["ocsp"] = func() (http.Handler, error) {
		src, err := ocsp_responder.NewSharedSources(ocspSigner)
		if err != nil {
			logger.Errorf("OCSP Sources Create error: %v", err)
			return nil, errors.Wrap(err, "sources Create error")
		}
		ocsp_responder.CountAll()
		return ocsp.NewResponder(src, nil), nil
	}

	registerHandlers()

	return router, nil
}

func tlsServe(addr string, tlsConfig *tls.Config) error {
	server := http.Server{
		Addr:      addr,
		TLSConfig: tlsConfig,
		Handler:   router,
	}
	return server.ListenAndServeTLS("", "")
}

// OcspServer
func OcspServer() ocsp.Signer {
	logger := core.Is.Logger.Named("singleca")
	ocspSigner, err := ocsp.NewDynamicSigner(
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
		}, 4*24*time.Hour)
	if err != nil {
		logger.Warnf("couldn't initialize ocsp signer: %v", err)
	}
	return ocspSigner
}
