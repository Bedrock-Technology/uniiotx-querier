package interactors

import (
	"github.com/Bedrock-Technology/uniiotx-querier/storer"
	"github.com/dgraph-io/ristretto"
)

type InteractorFactory struct {
	Cacher *ristretto.Cache
	Storer *storer.MyStorer
}
