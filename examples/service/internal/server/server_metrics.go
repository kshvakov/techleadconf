package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	info "github.com/kshvakov/techleadconf/examples/service"
)

var (
	promDBPool = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: info.Namespace,
		Name:      "db_pool",
		Help:      "The DB pool status.",
		ConstLabels: prometheus.Labels{
			"service": "example",
		},
	}, []string{"db", "status"})
)

func (srv *Server) PromHandler() func(http.ResponseWriter, *http.Request) {
	handler := promhttp.Handler()
	return func(w http.ResponseWriter, r *http.Request) {
		/*for db, stats := range map[string]sql.DBStats{
			"core": srv.conn.Stats(),
		} {
			promDBPool.WithLabelValues(db, "max_open_connections").Set(float64(stats.MaxOpenConnections))

			promDBPool.WithLabelValues(db, "idle").Set(float64(stats.Idle))
			promDBPool.WithLabelValues(db, "in_use").Set(float64(stats.InUse))
		}*/
		handler.ServeHTTP(w, r)
	}
}
