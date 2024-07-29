package utils

import (
	"math"
	"math/big"
)

func UniIOTXRewards(iotxRewards, exchangeRatio *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Mul(iotxRewards, big.NewInt(1e18))
	result.Div(result, exchangeRatio)
	return result
}

func BigIntToFloat64(number *big.Int, multiplier float64, precision int) float64 {
	val := new(big.Float).SetInt(number)
	val.Quo(val, big.NewFloat(multiplier))

	num, _ := val.Float64()
	mult := math.Pow10(precision)
	return math.Floor(num*mult) / mult
}
