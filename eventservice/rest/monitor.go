package rest

import (
	"github.com/prometheus/client_golang/prometheus"
)

var eventsAddCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "events_add_count",
		Namespace: "myevents",
		Help: "Amount of events created",
	},
)

var delayAddCount = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name: "events_add_count",
		Namespace: "myevents",
		Help: "Delay of events created",
		Buckets: []float64{10,50,100},	
	},
)

func init() {
	prometheus.MustRegister(eventsAddCount)
	prometheus.MustRegister(delayAddCount)
}
