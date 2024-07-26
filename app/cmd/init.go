package cmd

import (
	"github.com/Bedrock-Technology/uniiotx-querier/config"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&config.C.DevMode, "devMode", "", true, "Indicate whether to run in development mode")

	rootCmd.PersistentFlags().StringVarP(&config.C.LogFileName, "logFileName", "", "./log.log", "The file to which logs will be written. If left empty, logs will print to stderr and stdout")
	rootCmd.PersistentFlags().BoolVarP(&config.C.ConsoleEncoder, "consoleEncoder", "", false, "Indicate whether to log with console encoder")
	rootCmd.PersistentFlags().BoolVarP(&config.C.Stacktrace, "stacktrace", "", true, "Indicate whether to log with stacktrace")

	rootCmd.PersistentFlags().StringVarP(&config.C.DataServerAddr, "dataServerAddr", "", "0.0.0.0:8011", "Address to be used by data server")
	rootCmd.PersistentFlags().StringVarP(&config.C.MetricServerAddr, "metricServerAddr", "", "0.0.0.0:7000", "Address to be used by metric server")

	rootCmd.PersistentFlags().StringVarP(&config.C.SqliteDSN, "sqliteDSN", "", "./sqlite.db", "Sqlite data source name")

	rootCmd.PersistentFlags().StringVarP(&config.C.ChainHost, "chainHost", "", "https://babel-api.mainnet.iotex.io", "The blockchain host for RPC communication")

	rootCmd.PersistentFlags().StringVarP(&config.C.IOTXStaking, "iotxstaking", "", "0x2c914Ba874D94090Ba0E6F56790bb8Eb6D4C7e5f", "The address of IOTXStaking contract")
}
