package config

var C = &Config{
	DevMode: false,

	LogFileName:    "./log.log",
	ConsoleEncoder: false,
	Stacktrace:     true,

	DataServerAddr:   "0.0.0.0:8011",
	MetricServerAddr: "0.0.0.0:7000",

	SqliteDSN: "./sqlite.db",

	ChainHost: "https://babel-api.mainnet.iotex.io",

	SystemStaking: "0x68db92a6a78a39dcaff1745da9e89e230ef49d3d",
	IOTXStaking:   "0x2c914Ba874D94090Ba0E6F56790bb8Eb6D4C7e5f",
	IOTXClear:     "0x7AD800771743F4e29f55235A55895273035FB546",
}

type Config struct {
	DevMode bool

	LogFileName    string
	ConsoleEncoder bool
	Stacktrace     bool

	DataServerAddr   string
	MetricServerAddr string

	SqliteDSN string

	ChainHost string

	SystemStaking string
	IOTXStaking   string
	IOTXClear     string
}
