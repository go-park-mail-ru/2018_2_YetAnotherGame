package models

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct{
	Counter *prometheus.CounterVec
}