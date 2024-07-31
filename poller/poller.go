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
	ttl                           = time.Minute * 30
)

type Poller struct {
	Logger common.Logger

	SystemStakingCaller *bindings.SystemStakingCaller
	IOTXStakingCaller   *bindings.IOTXStakingCaller
	IOTXClearCaller     *bindings.IOTXClearCaller
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
			if err := p.syncAssetStatistics(); err != nil {
				continue
			}

			if err := p.syncBuckets(); err != nil {
				continue
			}
		}
	}
}

func (p *Poller) syncAssetStatistics() error {
	// 1. Query datum
	totalPending, err := p.IOTXStakingCaller.GetTotalPending(nil)
	if err != nil {
		p.Logger.Error("failed to query total pending value", err)
		return err
	}

	totalStaked, err := p.IOTXStakingCaller.GetTotalStaked(nil)
	if err != nil {
		p.Logger.Error("failed to query total staked value", err)
		return err
	}

	managerRewards, err := p.IOTXStakingCaller.GetManagerReward(nil)
	if err != nil {
		p.Logger.Error("failed to query manger rewards", err)
		return err
	}

	userRewards, err := p.IOTXStakingCaller.GetUserReward(nil)
	if err != nil {
		p.Logger.Error("failed to query user rewards", err)
		return err
	}

	exchangeRatio, err := p.IOTXStakingCaller.ExchangeRatio(nil)
	if err != nil {
		p.Logger.Error("failed to query exchange ratio", err)
		return err
	}

	totalDebts, err := p.IOTXClearCaller.TotalDebts(nil)
	if err != nil {
		p.Logger.Error("failed to query total debts", err)
		return err
	}

	date, year, month, day := utils.Date(time.Now())
	managerRewardsUniIOTX := utils.ConvertIOTXToUniIOTX(managerRewards, exchangeRatio)
	userRewardsUniIOTX := utils.ConvertIOTXToUniIOTX(userRewards, exchangeRatio)

	latestData := &common.DailyAssetStatistics{
		Date:  date,
		Year:  year,
		Month: month,
		Day:   day,

		TotalPending:  totalPending.String(),
		TotalStaked:   totalStaked.String(),
		TotalDebts:    totalDebts.String(),
		ExchangeRatio: exchangeRatio.String(),

		ManagerRewards:        managerRewards.String(),
		ManagerRewardsUniIOTX: managerRewardsUniIOTX.String(),
		UserRewards:           userRewards.String(),
		UserRewardsUniIOTX:    userRewardsUniIOTX.String(),
	}

	// 2. Update storage
	err = p.Storer.CreateOrUpdateDailyAssetStatistics(latestData, false)
	if err != nil {
		p.Logger.Error("failed to store manager rewards", err)
		return err
	}

	// 3. Update cache
	p.Cacher.SetWithTTL(common.CacheKeyLatestAssetStatistics, latestData, 1, ttl)
	p.Cacher.Wait()

	// 4. Update metrics
	totalPendingVal := utils.ConvertBigIntToFloat64(totalPending, 1e18, 3)
	totalStakedVal := utils.ConvertBigIntToFloat64(totalStaked, 1e18, 3)
	totalDebtsVal := utils.ConvertBigIntToFloat64(totalDebts, 1e18, 3)
	exchangeRatioVal := utils.ConvertBigIntToFloat64(exchangeRatio, 1e18, 3)

	managerRewardsVal := utils.ConvertBigIntToFloat64(managerRewards, 1e18, 3)
	managerRewardsUniIOTXVal := utils.ConvertBigIntToFloat64(managerRewardsUniIOTX, 1e18, 3)
	userRewardsVal := utils.ConvertBigIntToFloat64(userRewards, 1e18, 3)
	userRewardsUniIOTXVal := utils.ConvertBigIntToFloat64(userRewardsUniIOTX, 1e18, 3)

	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "totalPending"}).Set(totalPendingVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "totalStaked"}).Set(totalStakedVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "totalDebts"}).Set(totalDebtsVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "exchangeRatio"}).Set(exchangeRatioVal)

	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "managerRewards"}).Set(managerRewardsVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "managerRewardsUniIOTX"}).Set(managerRewardsUniIOTXVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "userRewards"}).Set(userRewardsVal)
	metrics.AssetStatistics.With(prometheus.Labels{"valueType": "userRewardsUniIOTX"}).Set(userRewardsUniIOTXVal)

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
