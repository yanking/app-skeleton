package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yanking/app-skeleton/pkg/log"
)

// A VectorOpts is a general configuration.
type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

func Handler() {
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("metrics: server listening on :9091")
		if err := http.ListenAndServe(":9091", metricsMux); err != nil {
			log.Fatalf("metrics: failed to start metrics server: %v", err)
		}
	}()
}
