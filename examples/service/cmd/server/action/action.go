package action

import (
	"github.com/kshvakov/techleadconf/examples/service/internal/server"
	"github.com/kshvakov/techleadconf/examples/service/pkg/intserver"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "http_addr",
			EnvVar: "HTTP_ADDR",
			Value:  "0.0.0.0:8000",
		},
	}
	Do = func(c *cli.Context) error {
		app, err := server.New(&server.Options{
			Addr: c.String("http_addr"),
		})
		if err != nil {
			return err
		}
		go func() {
			log.Infof("monitoring HTTP server listen on: %s", c.String("monitoring_http_addr"))
			if err := intserver.ListenAndServe(c.String("monitoring_http_addr"), app.PromHandler()); err != nil {
				log.Fatal("could not init monitoring server: ", err)
			}
		}()
		return app.Run()
	}
)
