package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	buildsTrigger = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pr_controller_builds_triggered",
		Help: "Count of CI Pipelines triggered",
	})

	environmentsActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pr_controller_environments_active",
		Help: "Count active pr environments",
	})

	errorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "pr_controller_errors",
		Help: "Count of errors",
	})
)

func Init() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":2112", nil)
}

func CIPipelineTriggered() {
	buildsTrigger.Inc()
}

func AddEnvironment() {
	environmentsActive.Inc()
}

func RemoveEnvironment() {
	environmentsActive.Dec()
}

func AddError() {
	errorCounter.Inc()
}
