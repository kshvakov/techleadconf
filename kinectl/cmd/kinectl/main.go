package main

import (
	"fmt"
	"os"

	"github.com/kshvakov/techleadconf/kinectl/cmd/kinectl/action"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/goreleaser/nfpm/v2/deb" // deb packager
	info "github.com/kshvakov/techleadconf/kinectl"
)

func init() {
	promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: info.Namespace,
		Name:      "application_info",
		Help:      "Information about the application",
		ConstLabels: prometheus.Labels{
			"name":         "kinectl",
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
		app.Usage = "kinectl"
		app.Version = fmt.Sprintf("%s (rev[%s] %s %s UTC).", info.AppVersion, info.GitCommit, info.GitBranch, info.ReleaseDate)
	}
	app.Commands = []cli.Command{
		action.Deb(),
		action.Play(),
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
