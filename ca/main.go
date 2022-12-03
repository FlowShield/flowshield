package main

import (
	"context"
	"github.com/flowshield/flowshield/ca/cmd"
	"github.com/flowshield/flowshield/ca/initer"
	"github.com/flowshield/flowshield/ca/pkg/logger"
	"github.com/urfave/cli"
	"os"
)

func main() {
	ctx := context.Background()

	app := cli.NewApp()
	app.Name = "capitalizone"
	app.Version = "1.0.0"
	app.Usage = "capitalizone"
	app.Commands = []cli.Command{
		newTlsCmd(ctx),
	}
	app.Before = initer.Init
	err := app.Run(os.Args)
	if err != nil {
		logger.Named("Init").Errorf(err.Error())
	}
}

// newTlsCmd Running TLS service
func newTlsCmd(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "tls",
		Usage: "Running TLS service",
		Action: func(c *cli.Context) error {
			return cmd.RunTls(ctx)
		},
	}
}
