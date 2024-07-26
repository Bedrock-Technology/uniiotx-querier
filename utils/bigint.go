package utils

import (
	"math/big"
)

func UniIOTXRewards(iotxRewards, exchangeRatio *big.Int) *big.Int {
	result := big.NewInt(0)
	result.Mul(iotxRewards, big.NewInt(1e18))
	result.Div(result, exchangeRatio)
	return result
}

func DivBy1e18(value *big.Int) *big.Int {
	value.Div(value, big.NewInt(1e18))
	return value
}

func DivBy1e14(value *big.Int) *big.Int {
	value.Div(value, big.NewInt(1e14))
	return value
}
