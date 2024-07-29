package poller

import (
	"github.com/Bedrock-Technology/uniiotx-querier/bindings"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/metrics"
	"github.com/Bedrock-Technology/uniiotx-querier/storer"
	"github.com/Bedrock-Technology/uniiotx-querier/utils"
	"github.com/dgraph-io/ristretto"
	"time"
)

const (
	pollingInterval = time.Second * 10
	ttl             = time.Minute * 15
)

type Poller struct {
	Logger common.Logger

	IOTXStakingCaller *bindings.IOTXStakingCaller
	Cacher            *ristretto.Cache
	Storer            *storer.MyStorer

	closeCh chan struct{}
}

func (p *Poller) Start() {
	p.closeCh = make(chan struct{})
	p.run()
}

func (p *Poller) Close() {
	close(p.closeCh)
}

func (p *Poller) run() {
	t := time.NewTicker(pollingInterval)
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
