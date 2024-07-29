package common

type DailyManagerRewards struct {
	Year           int    `json:"-" example:"2024" description:""`
	Month          int    `json:"-" example:"7" description:""`
	Day            int    `json:"-" example:"1" description:""`
	Date           int    `json:"date" example:"20240725" description:""`
	IOTXRewards    string `json:"iotxRewards" example:"668452753878420747506699" description:""`
	UniIOTXRewards string `json:"uniIotxRewards" example:"642747788815799551264093" description:""`
	ExchangeRatio  string `json:"exchangeRatio" example:"1039992304150870905" description:"uniIotxRewards = iotxRewards * 1e18 / exchangeRatio"`
}

type StakedDelegates struct {
	Total           int      `json:"total" example:"3" description:""`
	Level1Addresses []string `json:"level1Addresses" example:"[0x274880f6A49e272D014a38a6cBf70745F78be97c]" description:""`
	Level2Addresses []string `json:"level2Addresses" example:"[0x3F273789B560087bA93b6b23158D45b86ddBBDF4]" description:""`
	Level3Addresses []string `json:"level3Addresses" example:"[0xdc89A7A08EceA0105f509A992E2080d0dF90b5B8]" description:""`
}

type StakedBuckets struct {
	Delegate      string `json:"delegate" example:"0x274880f6A49e272D014a38a6cBf70745F78be97c" description:""`
	Total         int    `json:"total" example:"3" description:""`
	Level1Buckets []int  `json:"level1Buckets" example:"[28]" description:""`
	Level2Buckets []int  `json:"level2Buckets" example:"[30]" description:""`
	Level3Buckets []int  `json:"leve31Buckets" example:"[2038]" description:""`
}
