package main

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	info "github.com/kshvakov/techleadconf/examples/service"
	"github.com/kshvakov/techleadconf/examples/service/cmd/server/action"
)

func init() {
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: info.Namespace,
		Name:      "application_info",
		Help:      "Information about the application",
		ConstLabels: prometheus.Labels{
			"name":         "example-server",
			"branch":       info.GitBranch,
			"commit":       info.GitCommit,
			"version":      info.AppVersion,
			"release_date": info.ReleaseDate,
		},
	}).Add(1)
}

func main() {
	app := cli.NewApp()
	{
		app.Name = "Kinescope"
		app.Usage = "example"
		app.Version = fmt.Sprintf("%s (rev[%s] %s %s UTC).", info.AppVersion, info.GitCommit, info.GitBranch, info.ReleaseDate)
		app.Action = func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}
			return action.Do(c)
		}
	}
	app.Flags = append(action.Flags, []cli.Flag{
		cli.StringFlag{
			Name:   "monitoring_http_addr",
			EnvVar: "MONITORING_HTTP_ADDR",
			Value:  "0.0.0.0:9510",
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
		},
	}...)
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
