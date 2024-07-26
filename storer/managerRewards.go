package storer

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/storer/sqlc"
)

func (s *MyStorer) GetDailyManagerRewards(date int) (*common.DailyManagerRewards, error) {
	row, err := s.q.GetDailyManagerRewards(context.Background(), int64(date))
	if err != nil {
		return nil, err
	}

	return &common.DailyManagerRewards{
		Date:           date,
		Year:           int(row.Year),
		Month:          int(row.Month),
		Day:            int(row.Day),
		IOTXRewards:    row.Iotxrewards,
		UniIOTXRewards: row.Uniiotxrewards,
		ExchangeRatio:  row.Exchangeratio,
	}, nil
}

func (s *MyStorer) ListDailyManagerRewardsByYear(year int) ([]common.DailyManagerRewards, error) {
	rows, err := s.q.ListDailyManagerRewardsByYear(context.Background(), int64(year))
	if err != nil {
		return nil, err
	}

	var results []common.DailyManagerRewards
	for _, row := range rows {
		results = append(results, common.DailyManagerRewards{
			Date:           year,
			Year:           int(row.Year),
			Month:          int(row.Month),
			Day:            int(row.Day),
			IOTXRewards:    row.Iotxrewards,
			UniIOTXRewards: row.Uniiotxrewards,
			ExchangeRatio:  row.Exchangeratio,
		})
	}

	return results, nil
}

func (s *MyStorer) ListDailyManagerRewardsByMonth(year, month int) ([]common.DailyManagerRewards, error) {
	rows, err := s.q.ListDailyManagerRewardsByMonth(context.Background(), sqlc.ListDailyManagerRewardsByMonthParams{
		Year:  int64(year),
		Month: int64(month),
	})
	if err != nil {
		return nil, err
	}

	var results []common.DailyManagerRewards
	for _, row := range rows {
		results = append(results, common.DailyManagerRewards{
			Date:           year,
			Year:           int(row.Year),
			Month:          int(row.Month),
			Day:            int(row.Day),
			IOTXRewards:    row.Iotxrewards,
			UniIOTXRewards: row.Uniiotxrewards,
			ExchangeRatio:  row.Exchangeratio,
		})
	}

	return results, nil
}

func (s *MyStorer) CreateOrUpdateDailyManagerRewards(data *common.DailyManagerRewards, forceUpdate bool) error {
	date := int64(data.Date)
	row, err := s.q.GetDailyManagerRewards(context.Background(), date)
	//  Create a row if no existing record
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := s.q.CreateDailyManagerRewards(context.Background(), sqlc.CreateDailyManagerRewardsParams{
				Date:           date,
				Year:           int64(data.Year),
				Month:          int64(data.Month),
				Day:            int64(data.Day),
				Iotxrewards:    data.IOTXRewards,
				Uniiotxrewards: data.UniIOTXRewards,
				Exchangeratio:  data.ExchangeRatio,
			})

			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// No operation if no changes and not forced to update.
	if data.IOTXRewards == row.Iotxrewards && data.ExchangeRatio == row.Exchangeratio && !forceUpdate {
		return nil
	}

	// Update record if there are changes
	err = s.q.UpdateDailyManagerRewards(context.Background(), sqlc.UpdateDailyManagerRewardsParams{
		Date:           date,
		Iotxrewards:    data.IOTXRewards,
		Uniiotxrewards: data.UniIOTXRewards,
		Exchangeratio:  data.ExchangeRatio,
	})

	if err != nil {
		return err
	}
	return nil
}
