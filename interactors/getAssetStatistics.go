package interactors

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (i *InteractorFactory) GetAssetStatisticsFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct {
			Date int `json:"date" required:"true" example:"20240726" description:"The date to query asset statistics; 0 means the current day."`
		}

		type output struct {
			Data common.DailyAssetStatistics `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				// Read from cache for the latest record if the date is not specified
				if in.Date == 0 {
					var data *common.DailyAssetStatistics
					if val, ok := i.Cacher.Get(common.CacheKeyLatestAssetStatistics); ok {
						data = val.(*common.DailyAssetStatistics)
					}

					if data == nil {
						return status.Wrap(fmt.Errorf("failed to get manager rewards"), status.Unavailable)
					}

					out.Data = *data
					return nil
				}

				// Read from database if the date is specified
				data, err := i.Storer.GetDailyAssetStatistics(in.Date)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) {
						return status.Wrap(fmt.Errorf("failed to get manager rewards: %w", err), status.NotFound)
					}
					return status.Wrap(fmt.Errorf("failed to get manager rewards: %w", err), status.Internal)
				}

				out.Data = *data
				return nil
			})

		u.SetTitle("Get daily asset statistics")
		u.SetDescription("Get daily asset statistics. Data is synchronized from the blockchain every 15 seconds. " +
			"If the date input is 0, then the current day will be used instead.")
		u.SetExpectedErrors(status.Unavailable, status.Internal, status.NotFound)

		return u
	}
}
