package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/flowshield/flowshield/verifier/verify"

	"github.com/flowshield/flowshield/verifier/server"
	"github.com/urfave/cli"
	_ "go.uber.org/automaxprocs"
)

func main() {
	app := cli.NewApp()
	app.Name = "verifier"
	app.Author = "TS"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c",
			Value: "config.yaml",
			Usage: "config file url",
		},
	}
	app.Before = server.InitService
	app.Action = func(c *cli.Context) error {
		if err := verify.VerObj.Run(context.TODO()); err != nil {
			return err
		}
		server.RunHTTP()
		return nil
	}
	//app.After = func(c *cli.Context) error {
	//exitSignal()
	//return nil
	//}
	err := app.Run(os.Args)
	if err != nil {
		panic("app run error:" + err.Error())
	}
}

func exitSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for sig := range sigs {
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("Bye!")
			os.Exit(0)
		case syscall.SIGHUP:
			fmt.Println("+++++++++++++++++++++++++++++")
		default:
			fmt.Println(sig)
		}
	}
}
