package mtapp

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var M *Metrics

type Metrics struct {
	processLatency *prometheus.GaugeVec
}

func init() {
	prefix := "mtapp"
	M = &Metrics{
		processLatency: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name:      "process_latency",
				Namespace: prefix,
				Help:      "Время выполнения процесса",
			}, []string{
				"tid",
			}),
	}
}

func recordProcessLatency(tid ThreadID, duration time.Duration) {
	M.processLatency.WithLabelValues([]string{
		tid.String(),
	}...).Set(float64(duration.Milliseconds()))
}
