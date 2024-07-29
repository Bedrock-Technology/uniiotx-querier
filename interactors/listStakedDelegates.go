package interactors

import (
	"context"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/swaggest/usecase"
)

func (i *InteractorFactory) ListStakedDelegatesFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct{}

		type output struct {
			Data common.StakedDelegates `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				_ = in

				bucketLevelToDelegateList := make(map[int][]string)
				if val, ok := i.Cacher.Get(common.CacheKeyBucketLevelToDelegates); ok {
					bucketLevelToDelegateList = val.(map[int][]string)
				}

				delegateMap := make(map[string]struct{})
				for _, delegates := range bucketLevelToDelegateList {
					for _, delegate := range delegates {
						delegateMap[delegate] = struct{}{}
					}
				}

				delegateList := make([]string, 0)
				for delegate, _ := range delegateMap {
					delegateList = append(delegateList, delegate)
				}

				out.Data.Total = len(delegateList)
				out.Data.Level1Addresses = bucketLevelToDelegateList[0]
				out.Data.Level2Addresses = bucketLevelToDelegateList[1]
				out.Data.Level3Addresses = bucketLevelToDelegateList[2]

				return nil
			})

		u.SetTitle("List latest staked delegates")
		u.SetDescription("List latest staked delegates. Data is synchronized from the blockchain every 15 minutes.")

		return u
	}
}
