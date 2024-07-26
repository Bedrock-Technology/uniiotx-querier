package utils

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestDivBy1e18(t *testing.T) {
	val, _ := new(big.Int).SetString("10000000000000000000", 10)
	divVal := DivBy1e18(val)
	require.Equal(t, "10", divVal.String())
}
