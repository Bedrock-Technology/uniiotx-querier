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
