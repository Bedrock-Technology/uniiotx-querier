package interactors

import (
	"context"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/swaggest/usecase"
	"sort"
)

func (i *InteractorFactory) ListRedeemedBucketsFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct{}

		type output struct {
			Data common.RedeemedBuckets `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				_ = in

				buckets := make([]int, 0)
				if val, ok := i.Cacher.Get(common.CacheKeyRedeemedBuckets); ok {
					buckets = val.([]int)
				}

				out.Data.Total = len(buckets)
				out.Data.Buckets = buckets

				sort.Ints(out.Data.Buckets)

				return nil
			})

		u.SetTitle("List redeemed buckets")
		u.SetDescription("List redeemed buckets. Data is synchronized from the blockchain every 15 minutes.")

		return u
	}
}
