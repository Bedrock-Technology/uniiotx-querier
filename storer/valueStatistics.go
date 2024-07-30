package storer

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/Bedrock-Technology/uniiotx-querier/storer/sqlc"
)

func (s *MyStorer) GetDailyAssetStatistics(date int) (*common.DailyAssetStatistics, error) {
	row, err := s.q.GetDailyAssetStatistics(context.Background(), int64(date))
	if err != nil {
		return nil, err
	}

	return &common.DailyAssetStatistics{
		Date:  date,
		Year:  int(row.Year),
		Month: int(row.Month),
		Day:   int(row.Day),

		TotalPending:  row.Totalpending,
		TotalStaked:   row.Totalstaked,
		TotalDebts:    row.Totaldebts,
		ExchangeRatio: row.Exchangeratio,

		ManagerRewards:        row.Managerrewards,
		ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
		UserRewards:           row.Userrewards,
		UserRewardsUniIOTX:    row.Userrewardsuniiotx,
	}, nil
}

func (s *MyStorer) ListDailyAssetStatisticsByYear(year int) ([]common.DailyAssetStatistics, error) {
	rows, err := s.q.ListDailyAssetStatisticsByYear(context.Background(), int64(year))
	if err != nil {
		return nil, err
	}

	var results []common.DailyAssetStatistics
	for _, row := range rows {
		results = append(results, common.DailyAssetStatistics{
			Date:  int(row.Date),
			Year:  int(row.Year),
			Month: int(row.Month),
			Day:   int(row.Day),

			TotalPending:  row.Totalpending,
			TotalStaked:   row.Totalstaked,
			TotalDebts:    row.Totaldebts,
			ExchangeRatio: row.Exchangeratio,

			ManagerRewards:        row.Managerrewards,
			ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
			UserRewards:           row.Userrewards,
			UserRewardsUniIOTX:    row.Userrewardsuniiotx,
		})
	}

	return results, nil
}

func (s *MyStorer) ListDailyAssetStatisticsByMonth(year, month int) ([]common.DailyAssetStatistics, error) {
	rows, err := s.q.ListDailyAssetStatisticsByMonth(context.Background(), sqlc.ListDailyAssetStatisticsByMonthParams{
		Year:  int64(year),
		Month: int64(month),
	})
	if err != nil {
		return nil, err
	}

	var results []common.DailyAssetStatistics
	for _, row := range rows {
		results = append(results, common.DailyAssetStatistics{
			Date:  int(row.Date),
			Year:  int(row.Year),
			Month: int(row.Month),
			Day:   int(row.Day),

			TotalPending:  row.Totalpending,
			TotalStaked:   row.Totalstaked,
			TotalDebts:    row.Totaldebts,
			ExchangeRatio: row.Exchangeratio,

			ManagerRewards:        row.Managerrewards,
			ManagerRewardsUniIOTX: row.Managerrewardsuniiotx,
			UserRewards:           row.Userrewards,
			UserRewardsUniIOTX:    row.Userrewardsuniiotx,
		})
	}

	return results, nil
}

func (s *MyStorer) CreateOrUpdateDailyAssetStatistics(data *common.DailyAssetStatistics, forceUpdate bool) error {
	date := int64(data.Date)
	row, err := s.q.GetDailyAssetStatistics(context.Background(), date)
	//  Create a row if no existing record
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := s.q.CreateDailyAssetStatistics(context.Background(), sqlc.CreateDailyAssetStatisticsParams{
				Date:  date,
				Year:  int64(data.Year),
				Month: int64(data.Month),
				Day:   int64(data.Day),

				Totalpending:  data.TotalPending,
				Totalstaked:   data.TotalStaked,
				Totaldebts:    data.TotalDebts,
				Exchangeratio: data.ExchangeRatio,

				Managerrewards:        data.ManagerRewards,
				Managerrewardsuniiotx: data.ManagerRewardsUniIOTX,
				Userrewards:           data.UserRewards,
				Userrewardsuniiotx:    data.UserRewardsUniIOTX,
			})

			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// No operation if no changes and not forced to update.
	if data.TotalPending == row.Totalpending &&
		data.TotalStaked == row.Totalstaked &&
		data.TotalDebts == row.Totaldebts &&
		data.ExchangeRatio == row.Exchangeratio &&
		data.ManagerRewards == row.Managerrewards &&
		data.ManagerRewardsUniIOTX == row.Managerrewardsuniiotx &&
		data.UserRewards == row.Userrewards &&
		data.UserRewardsUniIOTX == row.Userrewardsuniiotx &&
		!forceUpdate {
		return nil
	}

	// Update record if there are changes
	err = s.q.UpdateDailyAssetStatistics(context.Background(), sqlc.UpdateDailyAssetStatisticsParams{
		Date: date,

		Totalpending:  data.TotalPending,
		Totalstaked:   data.TotalStaked,
		Totaldebts:    data.TotalDebts,
		Exchangeratio: data.ExchangeRatio,

		Managerrewards:        data.ManagerRewards,
		Managerrewardsuniiotx: data.ManagerRewardsUniIOTX,
		Userrewards:           data.UserRewards,
		Userrewardsuniiotx:    data.UserRewardsUniIOTX,
	})

	if err != nil {
		return err
	}
	return nil
}
