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

	TotalStakedDelegates = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: prefix + "total_staked_delegates",
			Help: "The total number of staked delegates",
		})

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
			Name: prefix + "asset_statistics",
			Help: "The statistics of asset info",
		}, []string{"valueType"})
)
