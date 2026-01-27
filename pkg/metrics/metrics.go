package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	MatchesMade = promauto.NewCounter(prometheus.CounterOpts{
		Name: "arena_matches_total",
		Help: "The total number of matches created",
	})

	PlayersInQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "arena_queue_depth",
		Help: "Current number of players waiting in Redis",
	})

	MatchLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "arena_match_latency_seconds",
		Help:    "Time taken to find a match",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
)
