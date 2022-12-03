package initer

import (
	"github.com/flowshield/flowshield/ca/ca/keymanager"
	"github.com/flowshield/flowshield/ca/core"
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/urfave/cli"
	"log"
	// ...
	_ "github.com/flowshield/flowshield/ca/util"
)

// Init Initialization
func Init(c *cli.Context) error {
	conf, err := parseConfigs(c)
	if err != nil {
		return err
	}
	initLogger(&conf)
	log.Printf("started with conf: %+v", conf)

	l := &core.Logger{Logger: logger.S()}

	i := &core.I{
		Config: &conf,
		Logger: l,
	}

	core.Is = i
	// CA Start
	if err := keymanager.InitKeeper(); err != nil {
		return err
	}

	logger.Info("success started.")
	return nil
}
