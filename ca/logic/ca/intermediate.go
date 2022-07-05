package ca

import (
	"crypto/tls"
	"net/http"

	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/helpers"
	"gorm.io/gorm"

	"github.com/cloudslit/cloudslit/ca/ca/upperca"
	"github.com/cloudslit/cloudslit/ca/core"
	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/dao"
	"github.com/cloudslit/cloudslit/ca/logic/schema"
	"github.com/cloudslit/cloudslit/ca/pkg/caclient"
)

const (
	UpperCaApiIntermediateTopology = "/api/v1/cap/ca/intermediate_topology"
)

var httpClient = resty.NewWithClient(&http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
})

type IntermediateObject struct {
	Certs    []*schema.FullCert    `mapstructure:"certs" json:"certs"`
	Metadata schema.CaMetadata     `mapstructure:"metadata" json:"metadata"`
	Children []*IntermediateObject `json:"children"`
	Current  bool                  `json:"current"`
}

// IntermediateTopology Obtain the sub cluster certificate issued by itself
func (l *Logic) IntermediateTopology() ([]*IntermediateObject, error) {
	db := l.db.Session(&gorm.Session{})
	db = db.Where("ca_label = ?", caclient.RoleIntermediate)
	db = db.Select(
		"ca_label",
		"issued_at",
		"serial_number",
		"authority_key_identifier",
		"status",
		"not_before",
		"expiry",
		"revoked_at",
		"pem",
	)
	list, _, err := dao.GetAllCertificates(db, 1, 100, "issued_at desc")
	if err != nil {
		return nil, errors.Wrap(err, "Database query error")
	}
	l.logger.Debugf("Number of query results: %v", len(list))
	intermediateMap := make(map[string]*IntermediateObject, 0)
	for _, row := range list {
		rawCert, err := helpers.ParseCertificatePEM([]byte(row.Pem))
		if err != nil {
			l.logger.With("row", row).Errorf("CA Certificate parsing error: %s", err)
			continue
		}
		if len(rawCert.Subject.OrganizationalUnit) == 0 || len(rawCert.Subject.Organization) == 0 {
			l.logger.With("row", row).Warn("CA Certificate missing O/OU Field")
			continue
		}
		ou := rawCert.Subject.OrganizationalUnit[0]
		if _, ok := intermediateMap[ou]; !ok {
			intermediateMap[ou] = &IntermediateObject{
				Metadata: schema.GetCaMetadataFromX509Cert(rawCert),
			}
		}
		intermediateMap[ou].Certs = append(intermediateMap[ou].Certs, schema.GetFullCertByX509Cert(rawCert))
	}

	result := make([]*IntermediateObject, 0, len(intermediateMap))
	for _, v := range intermediateMap {
		result = append(result, v)
	}

	return result, nil
}

// UpperCaIntermediateTopology Get parent CA's
func (l *Logic) UpperCaIntermediateTopology() ([]*IntermediateObject, error) {
	if core.Is.Config.Keymanager.SelfSign {
		return l.IntermediateTopology()
	}

	var resp *resty.Response
	err := upperca.ProxyRequest(func(host string) error {
		res, err := httpClient.R().Get(host + UpperCaApiIntermediateTopology)
		if err != nil {
			l.logger.With("upperca", host).Errorf("UpperCA Request error: %s", err)
			return err
		}
		resp = res
		return nil
	})
	if err != nil {
		l.logger.Errorf("UpperCA Sub CA topology acquisition failed: %s", err)
		return nil, err
	}

	body := resp.Body()
	var response struct {
		Data []*IntermediateObject `json:"data"`
	}
	if err := jsoniter.Unmarshal(body, &response); err != nil {
		l.logger.With("body", string(body)).Errorf("json Parsing error: %s", err)
		return nil, errors.Wrap(err, "json Parsing error")
	}

	return response.Data, nil
}
