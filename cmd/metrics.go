package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus Exporter Metrics
var (
	promProxyRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proxy_request_total",
		Help: "The total number of prroxy requests",
	})
)
