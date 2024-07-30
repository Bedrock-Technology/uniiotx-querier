package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	prefix = "uniiotx_"
)

var (
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

	AssetStatistics = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: prefix + "value_statistics",
			Help: "The statistics of value related info",
		}, []string{"valueType"})
)
