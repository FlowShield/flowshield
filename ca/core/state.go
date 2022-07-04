package core

import (
	"context"

	"github.com/cloudSlit/cloudslit/ca/core/config"
	"github.com/cloudSlit/cloudslit/ca/pkg/influxdb"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
	"github.com/cloudSlit/cloudslit/ca/pkg/vaultsecret"
	vaultAPI "github.com/hashicorp/vault/api"
	"gorm.io/gorm"
)

// Config ...
type Config struct {
	config.IConfig
}

// Is ...
var Is *I

// Elector ...
type Elector interface {
	IsLeader() bool
}

// Logger ...
type Logger struct {
	*logger.Logger
}

// I ...
type I struct {
	Ctx         context.Context
	Config      *Config
	Logger      *Logger
	Db          *gorm.DB
	Elector     Elector
	Metrics     *influxdb.Metrics
	VaultClient *vaultAPI.Client
	VaultSecret *vaultsecret.VaultSecret
}
