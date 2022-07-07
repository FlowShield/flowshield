package core

import (
	"context"

	"github.com/cloudslit/cloudslit/ca/core/config"
	"github.com/cloudslit/cloudslit/ca/pkg/logger"
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
	Ctx     context.Context
	Config  *Config
	Logger  *Logger
	Db      *gorm.DB
	Elector Elector
}
