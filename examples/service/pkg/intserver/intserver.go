package intserver

import (
	"net/http"
	"net/http/pprof"
)

func init() {
	http.DefaultServeMux = http.NewServeMux()
}

func Handle(mux *http.ServeMux, metrics http.HandlerFunc) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	if metrics != nil {
		mux.HandleFunc("/metrics", metrics)
	}
}

func NewServeMux(metrics http.HandlerFunc) *http.ServeMux {
	mux := http.NewServeMux()
	Handle(mux, metrics)
	return mux
}

func NewServer(addr string, metrics http.HandlerFunc) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: NewServeMux(metrics),
	}
}

func ListenAndServe(addr string, metrics http.HandlerFunc) error {
	return NewServer(addr, metrics).ListenAndServe()
}
