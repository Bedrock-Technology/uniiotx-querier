package storer

import (
	"context"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/storer/sqlc"
)

func (s *MyStorer) GetDailyManagerRewards(date int) (*common.DailyManagerRewards, error) {
	row, err := s.q.GetDailyAssetStatistics(context.Background(), int64(date))
	if err != nil {
		return nil, err
	}

	return &common.DailyManagerRewards{
		Date:  date,
		Year:  int(row.Year),
		Month: int(row.Month),
		Day:   int(row.Day),

		ExchangeRatio: row.Exchangeratio,

		ManagerRewards:        row.Managerrewards,
		ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
	}, nil
}

func (s *MyStorer) ListDailyManagerRewardsByYear(year int) ([]common.DailyManagerRewards, error) {
	rows, err := s.q.ListDailyAssetStatisticsByYear(context.Background(), int64(year))
	if err != nil {
		return nil, err
	}

	var results []common.DailyManagerRewards
	for _, row := range rows {
		results = append(results, common.DailyManagerRewards{
			Date:  int(row.Date),
			Year:  int(row.Year),
			Month: int(row.Month),
			Day:   int(row.Day),

			ExchangeRatio: row.Exchangeratio,

			ManagerRewards:        row.Managerrewards,
			ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
		})
	}

	return results, nil
}

func (s *MyStorer) ListDailyManagerRewardsByMonth(year, month int) ([]common.DailyManagerRewards, error) {
	rows, err := s.q.ListDailyAssetStatisticsByMonth(context.Background(), sqlc.ListDailyAssetStatisticsByMonthParams{
		Year:  int64(year),
		Month: int64(month),
	})
	if err != nil {
		return nil, err
	}

	var results []common.DailyManagerRewards
	for _, row := range rows {
		results = append(results, common.DailyManagerRewards{
			Date:  int(row.Date),
			Year:  int(row.Year),
			Month: int(row.Month),
			Day:   int(row.Day),

			ExchangeRatio: row.Exchangeratio,

			ManagerRewards:        row.Managerrewards,
			ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
		})
	}

	return results, nil
}
