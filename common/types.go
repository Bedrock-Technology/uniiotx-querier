package common

type DailyManagerRewards struct {
	Year           int    `json:"-" example:"2024" description:"The year when manager rewards are synchronized"`
	Month          int    `json:"-" example:"7" description:"The month when manager rewards are synchronized"`
	Day            int    `json:"-" example:"1" description:"The day when manager rewards are synchronized"`
	Date           int    `json:"date" example:"20240725" description:"The date when manager rewards are synchronized"`
	IOTXRewards    string `json:"iotxRewards" example:"668452753878420747506699" description:"The manager rewards in IOTX"`
	UniIOTXRewards string `json:"uniIotxRewards" example:"642747788815799551264093" description:"The manager rewards in uniIOTX"`
	ExchangeRatio  string `json:"exchangeRatio" example:"1039992304150870905" description:"uniIotxRewards = iotxRewards * 1e18 / exchangeRatio"`
}

type StakedDelegates struct {
	Total           int      `json:"total" example:"3" description:"The total number of delegates"`
	Level1Addresses []string `json:"level1Addresses" example:"[0x274880f6A49e272D014a38a6cBf70745F78be97c]" description:"Level 1 means the staking amount in its bucket is 10,000 IOTX"`
	Level2Addresses []string `json:"level2Addresses" example:"[0x3F273789B560087bA93b6b23158D45b86ddBBDF4]" description:"Level 2 means the staking amount in its bucket is 100,000 IOTX"`
	Level3Addresses []string `json:"level3Addresses" example:"[0xdc89A7A08EceA0105f509A992E2080d0dF90b5B8]" description:"Level 3 means the staking amount in its bucket is 1,000,000 IOTX"`
}

type StakedBuckets struct {
	Delegate      string `json:"delegate" example:"0x274880f6A49e272D014a38a6cBf70745F78be97c" description:"The address of the delegate"`
	Total         int    `json:"total" example:"3" description:"The total number of buckets"`
	Level1Buckets []int  `json:"level1Buckets" example:"[28]" description:"Level 1 means the staking amount in its bucket is 10,000 IOTX"`
	Level2Buckets []int  `json:"level2Buckets" example:"[30]" description:"Level 2 means the staking amount in its bucket is 100,000 IOTX"`
	Level3Buckets []int  `json:"leve31Buckets" example:"[2038]" description:"Level 3 means the staking amount in its bucket is 1,000,000 IOTX"`
}
