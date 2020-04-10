package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventSummary = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "events_processed",
			Help: "Events that have been processed",
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
	)
)

func init() {
	prometheus.MustRegister(eventSummary)
}

// EventCollector is used to collect the stats during event processing
func EventCollector(time int64) {
	// Just time for now
	eventSummary.Observe(float64(time))
}
