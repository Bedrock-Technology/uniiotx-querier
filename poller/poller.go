package poller

import (
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/bindings"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/metrics"
	"github.com/Bedrock-Technology/uniiotx-querier/storer"
	"github.com/Bedrock-Technology/uniiotx-querier/utils"
	"github.com/dgraph-io/ristretto"
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
	"time"
)

const (
	managerRewardsPollingInterval = time.Second * 15
	bucketsPollingInterval        = time.Minute * 15
	ttl                           = time.Minute * 15
)

type Poller struct {
	Logger common.Logger

	SystemStakingCaller *bindings.SystemStakingCaller
	IOTXStakingCaller   *bindings.IOTXStakingCaller
	Cacher              *ristretto.Cache
	Storer              *storer.MyStorer

	lastBucketSyncTime *time.Time

	closeCh  chan struct{}
	isClosed bool
}

func (p *Poller) Start() {
	p.closeCh = make(chan struct{})
	p.run()
}

func (p *Poller) Close() {
	close(p.closeCh)
	p.isClosed = true
}

func (p *Poller) run() {
	t := time.NewTicker(managerRewardsPollingInterval)
	defer t.Stop()

	p.Logger.Info("Poller started")
	for {
		select {
		case <-p.closeCh:
			p.Logger.Info("Poller closed")
			return
		case <-t.C:
			if err := p.syncManagerRewards(); err != nil {
				continue
			}

			if err := p.syncBuckets(); err != nil {
				continue
			}
		}
	}
}

func (p *Poller) syncManagerRewards() error {
	// 1. Query datum
	rewards, err := p.IOTXStakingCaller.GetManagerReward(nil)
	if err != nil {
		p.Logger.Error("failed to query manger rewards", err)
		return err
	}

	ratio, err := p.IOTXStakingCaller.ExchangeRatio(nil)
	if err != nil {
		p.Logger.Error("failed to query exchange ratio", err)
		return err
	}

	date, year, month, day := utils.Date(time.Now())
	uniIOTXRewards := utils.UniIOTXRewards(rewards, ratio)
	latestData := &common.DailyManagerRewards{
		Date:           date,
		Year:           year,
		Month:          month,
		Day:            day,
		IOTXRewards:    rewards.String(),
		UniIOTXRewards: uniIOTXRewards.String(),
		ExchangeRatio:  ratio.String(),
	}

	// 2. Update storage
	err = p.Storer.CreateOrUpdateDailyManagerRewards(latestData, false)
	if err != nil {
		p.Logger.Error("failed to store manager rewards", err)
		return err
	}

	// 3. Update cache
	var oldData *common.DailyManagerRewards
	if val, ok := p.Cacher.Get(common.CacheKeyLatestManagerRewards); ok {
		oldData = val.(*common.DailyManagerRewards)
	}

	p.Cacher.SetWithTTL(common.CacheKeyLatestManagerRewards, latestData, 1, ttl)
	p.Cacher.Wait()

	// 4. Update metrics
	rewardsVal := utils.BigIntToFloat64(rewards, 1e18, 3)
	uniIOTXRewardsVal := utils.BigIntToFloat64(uniIOTXRewards, 1e18, 3)
	ratioVal := utils.BigIntToFloat64(ratio, 1e18, 3)
	if oldData == nil {
		metrics.ManagerRewards.Set(rewardsVal)
		metrics.UniIOTXManagerRewards.Set(uniIOTXRewardsVal)
		metrics.ExchangeRatio.Set(ratioVal)
	} else {
		if latestData.IOTXRewards != oldData.IOTXRewards {
			metrics.ManagerRewards.Set(rewardsVal)
		}
		if latestData.UniIOTXRewards != oldData.UniIOTXRewards {
			metrics.UniIOTXManagerRewards.Set(uniIOTXRewardsVal)
		}
		if latestData.ExchangeRatio != oldData.ExchangeRatio {
			metrics.ExchangeRatio.Set(ratioVal)
		}
	}

	return nil
}

func (p *Poller) syncBuckets() error {
	// 1. Check time
	timeToSync := false

	if p.lastBucketSyncTime == nil || p.lastBucketSyncTime.Add(bucketsPollingInterval).Before(time.Now()) {
		timeToSync = true
	}

	if !timeToSync {
		return nil
	}

	// 2. Sync datum
	if err := p.syncStakedBuckets(); err != nil {
		return err
	}

	if err := p.syncRedeemedBuckets(); err != nil {
		return err
	}

	// 3. Update time
	now := time.Now()
	p.lastBucketSyncTime = &now
	return nil
}

func (p *Poller) syncStakedBuckets() error {
	// 1. Query datum
	seqLen, err := p.IOTXStakingCaller.SequenceLength(nil)
	if err != nil {
		p.Logger.Error("failed to query sequence length", err)
		return err
	}

	bucketLevelToBuckets := make(map[int][]int)

	for bucketLevel := 0; bucketLevel < int(seqLen.Int64()); bucketLevel++ {
		bucketBigs, err := p.IOTXStakingCaller.GetStakedTokenIds(nil, big.NewInt(int64(bucketLevel)))
		if err != nil {
			p.Logger.Error("failed to query token IDs", err)
			return err
		}
		for _, bucketBig := range bucketBigs {
			id := int(bucketBig.Int64())
			bucketLevelToBuckets[bucketLevel] = append(bucketLevelToBuckets[bucketLevel], id)
		}
	}

	delegateToBuckets := make(map[string]map[int][]int)
	bucketLevelToDelegates := make(map[int]map[string]struct{})

	// NOTE: The following algorithm performs poorly and can take too long
	for bucketLevel, buckets := range bucketLevelToBuckets {
		bucketLevelToDelegates[bucketLevel] = make(map[string]struct{})
		for _, bucket := range buckets {
			if p.isClosed {
				return nil
			}
			bucketStruct, err := p.SystemStakingCaller.BucketOf(nil, big.NewInt(int64(bucket)))
			if err != nil {
				p.Logger.Error("failed to query bucketStruct", err)
				return err
			}

			delegate := bucketStruct.Delegate.String()
			if _, ok := delegateToBuckets[delegate]; !ok {
				delegateToBuckets[delegate] = make(map[int][]int)
			}

			delegateToBuckets[delegate][bucketLevel] = append(delegateToBuckets[delegate][bucketLevel], bucket)
			bucketLevelToDelegates[bucketLevel][delegate] = struct{}{}
		}
	}

	// 2. Update cache
	bucketLevelToDelegateList := make(map[int][]string)
	for bucketLevel, delegates := range bucketLevelToDelegates {
		for delegate, _ := range delegates {
			bucketLevelToDelegateList[bucketLevel] = append(bucketLevelToDelegateList[bucketLevel], delegate)
		}
	}
	p.Cacher.SetWithTTL(common.CacheKeyStakedBucketLevelToDelegates, bucketLevelToDelegateList, 1, ttl)
	p.Cacher.Wait()

	p.Cacher.SetWithTTL(common.CacheKeyStakedDelegateToBuckets, delegateToBuckets, 1, ttl)
	p.Cacher.Wait()

	// 3. Update metrics
	for bucketLevel, buckets := range bucketLevelToBuckets {
		labelVal := fmt.Sprintf("bucketLevel%v", bucketLevel+1)
		metrics.StakedBuckets.With(prometheus.Labels{"bucketLevel": labelVal}).Set(float64(len(buckets)))
	}

	for bucketLevel, delegates := range bucketLevelToDelegateList {
		labelVal := fmt.Sprintf("bucketLevel%v", bucketLevel+1)
		metrics.StakedDelegates.With(prometheus.Labels{"bucketLevel": labelVal}).Set(float64(len(delegates)))
	}

	return nil
}

func (p *Poller) syncRedeemedBuckets() error {
	// 1. Query datum
	buckets := make([]int, 0)
	bucketBigs, err := p.IOTXStakingCaller.GetRedeemedTokenIds(nil)
	if err != nil {
		p.Logger.Error("failed to query token IDs", err)
		return err
	}
	for _, bucketBig := range bucketBigs {
		id := int(bucketBig.Int64())
		buckets = append(buckets, id)
	}

	// 2. Update cache
	p.Cacher.SetWithTTL(common.CacheKeyRedeemedBuckets, buckets, 1, ttl)
	p.Cacher.Wait()

	// 3. Update metrics
	metrics.RedeemedBuckets.Set(float64(len(buckets)))

	return nil
}
