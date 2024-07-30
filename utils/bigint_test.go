package utils

import (
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBigIntToFloat64(t *testing.T) {
	scenarios := []struct {
		bigIntStr  string
		float64Val float64
	}{
		{"678588874559603801060078", 678588.874},
		{"1040437090975081423", 1.04},
	}

	for _, s := range scenarios {
		val, _ := new(big.Int).SetString(s.bigIntStr, 10)
		result := ConvertBigIntToFloat64(val, 1e18, 3)
		require.Equal(t, s.float64Val, result)
	}
}
