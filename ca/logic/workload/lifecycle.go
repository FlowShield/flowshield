// Package workload Certificate lifecycle management
package workload

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ztalab/cfssl/ocsp"
	"gorm.io/gorm"

	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/dao"
	"github.com/cloudslit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudslit/cloudslit/ca/logic/events"
)

type RevokeCertsParams struct {
	SN       string `json:"sn"`
	AKI      string `json:"aki"`
	UniqueId string `json:"unique_id"`
}

// RevokeCerts Revocation of certificate
// 	1. Revoke certificate through snaki
//  2. Unified revocation of certificates through uniqueID
func (l *Logic) RevokeCerts(params *RevokeCertsParams) error {
	// 1. Certificate found by identity
	db := l.db.Session(&gorm.Session{})

	db = db.Where("status = ?", "good").
		Where("expiry > ?", time.Now())

	if params.UniqueId != "" {
		db = db.Where("common_name = ?", params.UniqueId)
	} else if params.AKI != "" && params.SN != "" {
		db = db.Where("serial_number = ? AND authority_key_identifier = ?", params.SN, params.AKI)
	} else {
		return errors.New("Parameter error")
	}

	certs, _, err := dao.GetAllCertificates(db, 1, 1000, "issued_at desc")
	if err != nil {
		l.logger.With("params", params).Errorf("Database query error: %s", err)
		return errors.Wrap(err, "Database query error")
	}

	if len(certs) == 0 {
		return errors.New("Certificate not found")
	}

	// 2. Batch revocation certificate
	reason, _ := ocsp.ReasonStringToCode("cacompromise")
	err = l.db.Transaction(func(tx *gorm.DB) error {
		for _, cert := range certs {
			err := tx.Model(&model.Certificates{}).Where(&model.Certificates{
				SerialNumber:           cert.SerialNumber,
				AuthorityKeyIdentifier: cert.AuthorityKeyIdentifier,
			}).Update("status", "revoked").
				Update("reason", reason).
				Update("revoked_at", time.Now()).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.logger.Errorf("Batch revocation certificate error: %s", err)
		return errors.Wrap(err, "Batch revocation certificate error")
	}

	// 3. Record operation log
	for _, cert := range certs {
		events.NewWorkloadLifeCycle("revoke", events.OperatorMSP, events.CertOp{
			UniqueId: cert.CommonName.String,
			SN:       cert.SerialNumber,
			AKI:      cert.AuthorityKeyIdentifier,
		}).Log()
	}

	return nil
}

type RecoverCertsParams struct {
	SN       string `json:"sn"`
	AKI      string `json:"aki"`
	UniqueId string `json:"unique_id"`
}

// RecoverCerts Restore certificate
// 	1. Recover certificate through snaki
//  2. Unified certificate recovery through uniqueID
func (l *Logic) RecoverCerts(params *RecoverCertsParams) error {
	// 1. Certificate found by identity
	db := l.db.Session(&gorm.Session{})

	db = db.Where("status = ?", "revoked").
		Where("expiry > ?", time.Now())

	switch {
	case params.UniqueId != "":
		db = db.Where("common_name = ?", params.UniqueId)
	case params.AKI != "" && params.SN != "":
		db = db.Where("serial_number = ? AND authority_key_identifier = ?", params.SN, params.AKI)
	default:
		return errors.New("Parameter error")
	}

	certs, _, err := dao.GetAllCertificates(db, 1, 1000, "issued_at desc")
	if err != nil {
		l.logger.With("params", params).Errorf("Database query error: %s", err)
		return errors.Wrap(err, "Database query error")
	}

	if len(certs) == 0 {
		return errors.New("Certificate not found")
	}

	// 2. Batch recovery certificate
	err = l.db.Transaction(func(tx *gorm.DB) error {
		for _, cert := range certs {
			err := tx.Model(&model.Certificates{}).Where(&model.Certificates{
				SerialNumber:           cert.SerialNumber,
				AuthorityKeyIdentifier: cert.AuthorityKeyIdentifier,
			}).Update("status", "good").
				Update("reason", 0).
				Update("revoked_at", nil).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 3. Record operation log
	for _, cert := range certs {
		events.NewWorkloadLifeCycle("recover", events.OperatorMSP, events.CertOp{
			UniqueId: cert.CommonName.String,
			SN:       cert.SerialNumber,
			AKI:      cert.AuthorityKeyIdentifier,
		}).Log()
	}

	return nil
}

type ForbidNewCertsParams struct {
	UniqueIds []string `json:"unique_ids"`
}

// ForbidNewCerts Prohibit a uniqueID from requesting a certificate
//	1.UniqueID is not allowed to apply for a new certificate
//  2. Logging
func (l *Logic) ForbidNewCerts(params *ForbidNewCertsParams) error {
	err := l.db.Transaction(func(tx *gorm.DB) error {
		for _, uid := range params.UniqueIds {
			record := model.Forbid{
				UniqueID:  uid,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, _, err := dao.AddForbid(tx, &record)
			if err != nil {
				l.logger.With("record", record).Errorf("Database insert error: %s", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.logger.Errorf("Database insert error: %s", err)
		return err
	}

	// Logging
	for _, uid := range params.UniqueIds {
		events.NewWorkloadLifeCycle("forbid", events.OperatorMSP, events.CertOp{
			UniqueId: uid,
		}).Log()
	}

	return nil
}

// RecoverForbidNewCerts Recovery allows a uniqueID to request a certificate
//	1. Allow uniqueID to request a new certificate
func (l *Logic) RecoverForbidNewCerts(params *ForbidNewCertsParams) error {
	err := l.db.Transaction(func(tx *gorm.DB) error {
		for _, uid := range params.UniqueIds {
			err := tx.Model(&model.Forbid{}).Where("unique_id = ?", uid).
				Where("deleted_at IS NULL").
				Update("deleted_at", time.Now()).Error
			if err != nil {
				l.logger.With("unique_id", uid).Errorf("Database update error: %s", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.logger.Errorf("Database update error: %s", err)
		return err
	}

	// Logging
	for _, uid := range params.UniqueIds {
		events.NewWorkloadLifeCycle("recover-forbid", events.OperatorMSP, events.CertOp{
			UniqueId: uid,
		}).Log()
	}

	return nil
}
