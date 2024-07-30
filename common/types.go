package common

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

type RedeemedBuckets struct {
	Total   int   `json:"total" example:"3" description:"The total number of buckets"`
	Buckets []int `json:"Buckets" example:"[2038]" description:"The bucket ID list with an amount of 1,000,000 IOTX"`
}

type DailyAssetStatistics struct {
	Date  int `json:"date" example:"20240725" description:"The date when manager rewards are synchronized"`
	Year  int `json:"-" example:"2024" description:"The year when statistics are synchronized"`
	Month int `json:"-" example:"7" description:"The month when statistics are synchronized"`
	Day   int `json:"-" example:"1" description:"The day when statistics are synchronized"`

	TotalPending  string `json:"totalPending" example:"136546761896526753742786" description:"The total pending in IOTX"`
	TotalStaked   string `json:"totalStaked" example:"450430000000000000000000000" description:"The total staked in IOTX"`
	TotalDebts    string `json:"totalDebts" example:"0" description:"The total debts in IOTX"`
	ExchangeRatio string `json:"exchangeRatio" example:"1040589943347730138" description:"uniIotxRewards = iotxRewards * 1e18 / exchangeRatio"`

	ManagerRewards        string `json:"managerRewards" example:"682072230887221328528316" description:"The manager rewards in IOTX"`
	ManagerRewardsUniIOTX string `json:"managerRewardsUniIOTX" example:"655466867854685526318582" description:"The manager rewards in uniIOTX"`
	UserRewards           string `json:"userRewards" example:"12959372386857205242040281" description:"The user rewards in IOTX"`
	UserRewardsUniIOTX    string `json:"userRewardsUniIOTX" example:"12453870489239025000055249" description:"The user rewards in uniIOTX"`
}

type DailyManagerRewards struct {
	Date  int `json:"date" example:"20240725" description:"The date when manager rewards are synchronized"`
	Year  int `json:"-" example:"2024" description:"The year when manager rewards are synchronized"`
	Month int `json:"-" example:"7" description:"The month when manager rewards are synchronized"`
	Day   int `json:"-" example:"1" description:"The day when manager rewards are synchronized"`

	ExchangeRatio string `json:"exchangeRatio" example:"1040589943347730138" description:"uniIotxRewards = iotxRewards * 1e18 / exchangeRatio"`

	ManagerRewards        string `json:"managerRewards" example:"682072230887221328528316" description:"The manager rewards in IOTX"`
	ManagerRewardsUniIOTX string `json:"managerRewardsUniIOTX" example:"655466867854685526318582" description:"The manager rewards in uniIOTX"`
}
