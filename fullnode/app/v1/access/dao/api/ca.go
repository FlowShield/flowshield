package api

import (
	"crypto/x509/pkix"
	"encoding/hex"
	"os"
	"sync"
	"time"

	"github.com/flowshield/casdk/caclient"
	"github.com/flowshield/casdk/keygen"
	"github.com/flowshield/casdk/pkg/attrmgr"
	"github.com/flowshield/casdk/pkg/spiffe"
	"github.com/flowshield/cfssl/helpers"
	"github.com/flowshield/flowshield/fullnode/pkg/confer"
	"github.com/flowshield/flowshield/fullnode/pkg/logger"
	"github.com/gin-gonic/gin"
)

var once sync.Once
var caClient *caclient.CAInstance

type SentinelSign struct {
	CaPEM     string
	CertPEM   string
	KeyPEM    string
	Sn        string
	Aki       string
	ExpiredAt time.Time
}

func getCaClient() *caclient.CAInstance {
	once.Do(func() {
		if caClient == nil {
			cfg := confer.GlobalConfig().CA
			caClient = caclient.NewCAI(
				caclient.WithCAServer(caclient.RoleDefault /*哨兵*/, cfg.SignURL),
				caclient.WithAuthKey(cfg.AuthKey),
			)
		}
	})
	return caClient
}

func ApplySign(c *gin.Context, attrs map[string]interface{}, uniqueID, cn, host string, duration time.Duration) (sentinelSign SentinelSign, err error) {
	client := getCaClient()
	mgr, err := client.NewCertManager()
	if err != nil {
		return
	}
	// CA PEM
	caPEMBytes, err := mgr.CACertsPEM()
	if err != nil {
		logger.Errorf(c, "mgr.CACertsPEM() err : %v", err)
		return
	}
	caPEM := string(caPEMBytes)
	// KEY PEM
	_, keyPEMBytes, _ := keygen.GenKey(keygen.EcdsaSigAlg)
	// 证书扩展字段
	attr := attrmgr.New()
	ext, _ := attr.ToPkixExtension(&attrmgr.Attributes{
		// 注入参数 Map[string]interface{}
		Attrs: attrs,
	})
	// gen csr
	csrPEM, _ := keygen.GenCustomExtendCSR(keyPEMBytes, &spiffe.IDGIdentity{
		SiteID:    os.Getenv("SITEUID"), /* Site 标识 */
		ClusterID: os.Getenv("CLUSTERUID"),
		UniqueID:  uniqueID,
	}, &keygen.CertOptions{ /* 通常为固定值 */
		CN:   cn,
		Host: host,
		TTL:  duration,
	}, []pkix.Extension{ext} /* 注入扩展字段 */)
	// get cert
	certPEMBytes, err := mgr.SignPEM(csrPEM, nil)
	if err != nil {
		logger.Errorf(c, "mgr.SignPEM() err : %v", err)
		return
	}
	certPEM := string(certPEMBytes)
	cert, err := helpers.ParseCertificatePEM(certPEMBytes)
	if err != nil {
		logger.Errorf(c, "helpers.ParseCertificatePEM() err : %v", err)
		return
	}
	sentinelSign = SentinelSign{
		CaPEM:   caPEM,
		CertPEM: certPEM,
		KeyPEM:  string(keyPEMBytes),
		Sn:      cert.SerialNumber.String(),
		Aki:     hex.EncodeToString(cert.AuthorityKeyId),
		//ExpiredAt: expiredAt,
		//ExpiredAt: cert.NotAfter,
	}
	return
}
