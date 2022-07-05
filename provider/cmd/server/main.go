package main

import (
	"context"
	cli1 "github.com/cloudslit/cloudslit/provider/cmd/server/cli"
	"os"

	"github.com/cloudslit/cloudslit/provider/pkg/logger"
	"github.com/urfave/cli/v2"
)

var VERSION = "0.0.0"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTagContext(context.Background(), "__main__")

	app := cli.NewApp()
	app.Name = "za-sentinel"
	app.Version = VERSION
	app.Usage = "Security, network acceleration, zero trust network architecture"
	app.Flags = commonConfig()
	app.Commands = cli1.NewServerCmd(ctx)
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf("%v", err)
	}
}

func commonConfig() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Usage:   "App configuration file(.json,.yaml,.toml)",
			Value:   "configs/config_server.toml",
		},
	}
}
