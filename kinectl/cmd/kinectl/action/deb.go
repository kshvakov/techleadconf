package action

import (
	"github.com/kshvakov/techleadconf/kinectl/internal/command/deb"
	"github.com/urfave/cli"
)

func Deb() cli.Command {
	return cli.Command{
		Name:  "deb",
		Usage: "creates a DEB package",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "spec",
				Value: "spec.yml",
			},
			cli.StringFlag{
				Name:  "target",
				Value: ".build",
			},
			cli.StringFlag{
				Name:     "app_version",
				EnvVar:   "APP_VERSION",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			app, err := deb.New(&deb.Options{
				Target:   c.String("target"),
				Version:  c.String("app_version"),
				SpecFile: c.String("spec"),
			})
			if err != nil {
				return err
			}
			return app.Run()
		},
	}
}
