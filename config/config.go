package config

var C = &Config{
	DevMode: true,

	LogFileName:    "./log.log",
	ConsoleEncoder: false,
	Stacktrace:     true,

	DataServerAddr:   "0.0.0.0:8011",
	MetricServerAddr: "0.0.0.0:7000",

	SqliteDSN: "./sqlite.db",

	ChainHost: "https://babel-api.mainnet.iotex.io",

	IOTXStaking: "0x2c914Ba874D94090Ba0E6F56790bb8Eb6D4C7e5f",
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

	IOTXStaking string
}
