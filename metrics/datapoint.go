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
			Help: "The amount of manager rewards of IOTX",
		})

	UniIOTXManagerRewards = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "manager_rewards_uniiotx",
			Help: "The amount of manager rewards of uniIOTX",
		})

	ExchangeRatio = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "exchange_ratio",
			Help: "The exchange ratio of uniIOTX to IOTX",
		})

	StakedDelegates = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "staked_delegates",
			Help: "The number of staked delegates at specified bucket level",
		}, []string{"bucketLevel"})

	StakedBuckets = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "staked_buckets",
			Help: "The number of staked buckets at specified bucket level",
		}, []string{"bucketLevel"})

	RedeemedBuckets = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "redeemed_buckets",
			Help: "The number of redeemed buckets",
		})
)
