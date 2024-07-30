package interactors

import (
	"context"
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"time"
)

func (i *InteractorFactory) ListAssetStatisticsByYearFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct {
			Year int `json:"year" required:"true" example:"2024" description:"The year to query asset statistics; 0 means the current year"`
		}

		type output struct {
			Data []common.DailyAssetStatistics `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				// Use current year instead if the year is not specified
				year := in.Year
				currentY, _, _ := time.Now().Date()
				if year == 0 {
					year = currentY
				}
				data, err := i.Storer.ListDailyAssetStatisticsByYear(year)
				if err != nil {
					return status.Wrap(fmt.Errorf("failed to list manager rewards: %w", err), status.Internal)
				}

				out.Data = data
				return nil
			})

		u.SetTitle("List yearly asset statistics")
		u.SetDescription("List yearly asset statistics. Data is synchronized from the blockchain every 15 seconds. " +
			"If the year input is 0, then the current year will be used instead.")
		u.SetExpectedErrors(status.Internal)

		return u
	}
}

func (i *InteractorFactory) ListAssetStatisticsByMonthFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct {
			Year  int `json:"year" required:"true" example:"2024" description:"The year to query asset statistics; 0 means the current year"`
			Month int `json:"month" required:"true" example:"7" description:"The month to query asset statistics; 0 means the current month"`
		}

		type output struct {
			Data []common.DailyAssetStatistics `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				// Use current year && month instead if the year && month is not specified
				year := in.Year
				month := in.Month
				currentY, currentM, _ := time.Now().Date()
				if year == 0 {
					year = currentY
				}
				if month == 0 {
					month = int(currentM)
				}

				data, err := i.Storer.ListDailyAssetStatisticsByMonth(year, month)
				if err != nil {
					return status.Wrap(fmt.Errorf("failed to list manager rewards: %w", err), status.Internal)
				}

				out.Data = data
				return nil
			})

		u.SetTitle("List monthly asset statistics")
		u.SetDescription("List monthly asset statistics. Data is synchronized from the blockchain every 15 seconds. " + "" +
			"If the year and month inputs are 0, then the current year and current month will be used instead.")
		u.SetExpectedErrors(status.Internal)

		return u
	}
}
