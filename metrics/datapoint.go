package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	prefix = "uniiotx_"
)

var (
	ManagerRewards = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "manager_rewards",
			Help: "The manager rewards of IOTX",
		})

	UniIOTXManagerRewards = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "manager_rewards_uniiotx",
			Help: "The manager rewards of uniIOTX",
		})

	ExchangeRatio = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "exchange_ratio",
			Help: "The exchange ratio of uniIOTX to IOTX",
		})

	Delegates = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "delegates",
			Help: "The latest delegates of specified bucket amount level",
		}, []string{"bucketLevel"})

	Buckets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "buckets",
			Help: "The latest buckets of specified bucket amount level",
		}, []string{"bucketLevel"})
)
