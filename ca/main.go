package main

import (
	"context"
	"github.com/cloudSlit/cloudslit/ca/cmd"
	"github.com/cloudSlit/cloudslit/ca/initer"
	"github.com/cloudSlit/cloudslit/ca/pkg/logger"
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
		newApiCmd(ctx),
		newTlsCmd(ctx),
		newOcspCmd(ctx),
	}
	app.Before = initer.Init
	err := app.Run(os.Args)
	if err != nil {
		logger.Named("Init").Errorf(err.Error())
	}
}

// newApiCmd Running API services
func newApiCmd(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "api",
		Usage: "Running API service",
		Action: func(c *cli.Context) error {
			return cmd.RunHttp(ctx)
		},
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

// newOcspCmd Running OCSP service
func newOcspCmd(ctx context.Context) cli.Command {
	return cli.Command{
		Name:  "ocsp",
		Usage: "Run OCSP service",
		Action: func(c *cli.Context) error {
			return cmd.RunOcsp(ctx)
		},
	}
}
