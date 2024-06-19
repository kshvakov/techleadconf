package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Options struct {
	Addr string
}

func New(o *Options) (*Server, error) {
	return &Server{
		opt: o,
	}, nil
}

type Server struct {
	opt *Options
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, "home")
	case "/info":
		fmt.Fprint(w, "info")
	default:
		http.NotFound(w, r)
	}
}

func (srv *Server) Run() error {
	server := http.Server{
		Addr:    srv.opt.Addr,
		Handler: srv,
	}
	log.Infof("example    HTTP server listen on: %s", srv.opt.Addr)
	return server.ListenAndServe()
}
