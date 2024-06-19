package action

import (
	"github.com/kshvakov/techleadconf/kinectl/internal/command/play"
	"github.com/urfave/cli"
)

func Play() cli.Command {
	return cli.Command{
		Name:  "play",
		Usage: "creates and apply an Ansible playbooks",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:     "env",
				Usage:    "(dev, staging, production etc...)",
				Required: true,
			},
			cli.StringFlag{
				Name:  "spec",
				Value: "spec.yml",
			},
			cli.StringFlag{
				Name:  "ops-dir",
				Value: "ops",
			},
			cli.StringFlag{
				Name: "limit",
			},
			cli.BoolFlag{
				Name: "dry-run",
			},
			cli.BoolFlag{
				Name: "diff",
			},
			cli.BoolFlag{
				Name: "debug",
			},
		},
		Action: func(c *cli.Context) error {
			app, err := play.New(&play.Options{
				Env:      c.String("env"),
				Diff:     c.Bool("diff"),
				Debug:    c.Bool("debug"),
				Limit:    c.String("limit"),
				DryRun:   c.Bool("dry-run"),
				OpsDir:   c.String("ops-dir"),
				SpecFile: c.String("spec"),
			})
			if err != nil {
				return err
			}
			return app.Run()
		},
	}
}
