package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/cloudslit/cloudslit/provider/internal"
	"github.com/cloudslit/cloudslit/provider/pkg/logger"
)

var VERSION = "0.0.0"

func main() {
	logger.SetVersion(VERSION)
	ctx := logger.NewTagContext(context.Background(), "__main__")

	app := cli.NewApp()
	app.Name = "provider"
	app.Version = VERSION
	app.Usage = "Security, network acceleration, zero trust network architecture"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "conf",
			Aliases: []string{"c"},
			Usage:   "App configuration file(.toml .json .yaml)",
			Value:   "configs/config.toml",
		},
	}
	app.Action = func(c *cli.Context) error {
		return internal.Run(ctx,
			internal.SetConfigFile(c.String("conf")))
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf("%v", err)
	}
}
