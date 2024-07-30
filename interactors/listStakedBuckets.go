package interactors

import (
	"context"
	"github.com/Bedrock-Technology/uniiotx-querier/common"
	"github.com/swaggest/usecase"
	"sort"
)

func (i *InteractorFactory) ListStakedBucketsFn() func() usecase.IOInteractor {
	return func() usecase.IOInteractor {
		type input struct {
			Delegate string `json:"delegate" example:"0xe535379690Dc22Dec14C999E263951C4127B5ACD" description:"The delegate for listing staked buckets"`
		}

		type output struct {
			Data common.StakedBuckets `json:"data"`
		}

		u := usecase.NewIOI(new(input), new(output),
			func(ctx context.Context, _input, _output any) (err error) {
				var (
					in  = _input.(*input)
					out = _output.(*output)
				)

				_ = in

				delegateToBuckets := make(map[string]map[int][]int)
				if val, ok := i.Cacher.Get(common.CacheKeyStakedDelegateToBuckets); ok {
					delegateToBuckets = val.(map[string]map[int][]int)
				}

				bucketLevelToBuckets := make(map[int][]int)
				if in.Delegate == "" {
					for _, bucketLevelToBucketsMap := range delegateToBuckets {
						for bucketLevel, buckets := range bucketLevelToBucketsMap {
							bucketLevelToBuckets[bucketLevel] = append(bucketLevelToBuckets[bucketLevel], buckets...)
						}
					}
				} else {
					bucketLevelToBuckets, _ = delegateToBuckets[in.Delegate]
				}

				out.Data.Delegate = in.Delegate
				if len(bucketLevelToBuckets) > 0 {
					sort.Ints(bucketLevelToBuckets[0])
					sort.Ints(bucketLevelToBuckets[1])
					sort.Ints(bucketLevelToBuckets[2])

					out.Data.Total = len(bucketLevelToBuckets[0]) + len(bucketLevelToBuckets[1]) + len(bucketLevelToBuckets[2])
					out.Data.Level1Buckets = bucketLevelToBuckets[0]
					out.Data.Level2Buckets = bucketLevelToBuckets[1]
					out.Data.Level3Buckets = bucketLevelToBuckets[2]
				}

				return nil
			})

		u.SetTitle("List the latest staked buckets")
		u.SetDescription("List the latest staked buckets. Data is synchronized from the blockchain every 15 minutes. " +
			"If the delegate input is left empty, then all staked buckets will be taken into account.")

		return u
	}
}
