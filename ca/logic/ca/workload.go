// Package ca workload Relevant
package ca

import (
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/tal-tech/go-zero/core/fx"
	"gorm.io/gorm"

	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/dao"
	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudSlit/cloudslit/ca/pkg/caclient"
	"github.com/cloudSlit/cloudslit/ca/util"
)

const AllCertsCacheKey = "all_certs_cache"

// WorkloadUnit UniqueID Divided workload unit
type WorkloadUnit struct {
	Role          caclient.Role `json:"role"`
	ValidNum      int           `json:"valid_num"`       // Number of valid certificates
	FirstIssuedAt time.Time     `json:"first_issued_at"` // Date of first issuance of certificate
	UniqueId      string        `json:"unique_id"`
	Forbidden     bool          `json:"forbidden"` // Is it prohibited
}

type WorkloadUnitsParams struct {
	Page, PageSize int
	UniqueId       string
}

// WorkloadUnits CA Units
// Return to currently active units and summary
func (l *Logic) WorkloadUnits(params *WorkloadUnitsParams) ([]*WorkloadUnit, int64, error) {
	db := l.db.Session(&gorm.Session{})
	// The default filter has no expired
	db = db.Where("expiry > ?", time.Now()).
		Where("status", "good")
	db = db.Select(
		"ca_label",
		"common_name",
		"issued_at",
		"serial_number",
		"authority_key_identifier",
		"status",
		"not_before",
		"expiry",
		"revoked_at",
	)

	certs, err := getCerts(db)
	if err != nil {
		return make([]*WorkloadUnit, 0), 0, errors.Wrap(err, "Database query error")
	}

	var i int
	var total int64
	units := make([]*WorkloadUnit, 0)
	fx.From(func(source chan<- interface{}) {
		for _, cert := range certs {
			source <- cert
		}
	}).Group(func(item interface{}) interface{} {
		cert := item.(*model.Certificates)
		return cert.CommonName
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		certs := item.([]interface{})
		firstCert := certs[0].(*model.Certificates)
		for _, certObj := range certs {
			cert := certObj.(*model.Certificates)
			if cert.IssuedAt.Before(firstCert.IssuedAt) {
				firstCert = cert
			}
		}
		unit := &WorkloadUnit{
			Role:          caclient.Role(firstCert.CaLabel.String),
			ValidNum:      len(certs),
			FirstIssuedAt: firstCert.IssuedAt,
			UniqueId:      firstCert.CommonName.String,
		}
		pipe <- unit
		atomic.AddInt64(&total, 1)
	}).Filter(func(item interface{}) bool {
		unit := item.(*WorkloadUnit)
		if params.UniqueId != "" {
			return unit.UniqueId == params.UniqueId
		}
		return true
	}).Sort(func(a, b interface{}) bool {
		aObj := a.(*WorkloadUnit)
		bObj := b.(*WorkloadUnit)
		return aObj.FirstIssuedAt.Before(bObj.FirstIssuedAt)
	}).Split(params.PageSize).ForEach(func(item interface{}) {
		i++
		if i == params.Page {
			group := item.([]interface{})
			for _, obj := range group {
				unit := obj.(*WorkloadUnit)
				units = append(units, unit)
			}
		}
	})

	return units, total, nil
}

func getCerts(db *gorm.DB) ([]*model.Certificates, error) {
	var certs []*model.Certificates
	var err error
	allCerts, ok := util.MapCache.Get(AllCertsCacheKey)
	if !ok {
		certs, _, err = dao.GetAllCertificates(db, 1, 10000, "issued_at desc")
		if err != nil {
			return nil, errors.Wrap(err, "Database query error")
		}
		util.MapCache.SetDefault(AllCertsCacheKey, certs)
	}
	if allCerts != nil {
		certs = allCerts.([]*model.Certificates)
	}
	return certs, nil
}
