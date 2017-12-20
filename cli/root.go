package cli

import (
	"context"
	"time"

	"github.com/corpix/loggers"
	"github.com/urfave/cli"

	"github.com/cryptounicorns/gluttony/input"
)

var (
	// RootCommands is a list of subcommands for the application.
	RootCommands = []cli.Command{}

	// RootFlags is a list of flags for the application.
	RootFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "application configuration file",
			EnvVar: "CONFIG",
			Value:  "config.json",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "add this flag to enable debug mode",
		},
	}
)

func runInput(c input.Config, l loggers.Logger) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		i      *input.Input
		err    error
	)

	for {
		// FIXME: What a mess...
		if cancel != nil {
			cancel()
		}
		if i != nil {
			i.Close()
		}
		if err != nil {
			time.Sleep(5 * time.Second)
		}

		i, err = input.New(c, l)
		if err != nil {
			l.Error(err)
			continue
		}

		ctx, cancel = context.WithCancel(context.Background())

		err = i.Run(ctx)
		if err != nil {
			l.Error(err)
			continue
		}
	}
}

// RootAction is executing when program called without any subcommand.
func RootAction(c *cli.Context) error {
	for _, i := range Config.Inputs {
		log.Printf(
			"Running '%s'",
			i.Name,
		)
		go runInput(i, log)
	}

	select {}

	return nil
}
