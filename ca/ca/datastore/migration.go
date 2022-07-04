package datastore

import (
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"time"

	"github.com/cloudSlit/cloudslit/ca/core"
	"github.com/cloudSlit/cloudslit/ca/database/mysql/cfssl-model/model"
	"github.com/cloudSlit/cloudslit/ca/pkg/vaultsecret"
)

// RunMigration Migrate MySQL data to vault
func RunMigration() {
	logger.Debug("MySQL -> Vault Database migration")
	certRows := make([]*model.Certificates, 0)
	result := core.Is.Db.Model(&model.Certificates{}).Where("expiry > ? AND revoked_at is NULL", time.Now()).Limit(10000).
		Find(&certRows)
	for _, row := range certRows {
		if row.Pem == "" {
			continue
		}
		if pemStr, err := core.Is.VaultSecret.GetCertPEM(row.SerialNumber); err != nil || *pemStr == "" {
			core.Is.Logger.Debugf("Vault Transfer %s", row.SerialNumber)
			if err := core.Is.VaultSecret.StoreCertPEM(row.SerialNumber, row.Pem); err != nil {
				core.Is.Logger.Errorf("Vault store cert %s Error: %s", row.SerialNumber, err)
			}
		}
	}
	if result.Error != nil {
		core.Is.Logger.Errorf("Error migrating Mysql to vault: %s", result.Error)
	}

	caKeyPair := new(model.SelfKeypair)
	if err := core.Is.Db.Model(&model.SelfKeypair{}).Where("name = ?", "ca").First(caKeyPair).Error; err == nil {
		if err := core.Is.VaultSecret.StoreCertPEMKey(vaultsecret.CALocalStoreKey,
			caKeyPair.Certificate.String, caKeyPair.PrivateKey.String); err != nil {
			core.Is.Logger.Errorf("Vault ca cert Storage error: %s", err)
		}
	}

	trustKeyPair := new(model.SelfKeypair)
	if err := core.Is.Db.Model(&model.SelfKeypair{}).Where("name = ?", "trust").First(trustKeyPair).Error; err == nil {
		if err := core.Is.VaultSecret.StoreCertPEM(vaultsecret.CATructCertsKey, trustKeyPair.Certificate.String); err != nil {
			core.Is.Logger.Errorf("Vault trust cert Storage error: %s", err)
		}
	}
}
