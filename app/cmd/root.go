package cmd

import (
	"context"
	"fmt"
	"github.com/Bedrock-Technology/uniiotx-querier/bindings"
	"github.com/Bedrock-Technology/uniiotx-querier/config"
	"github.com/Bedrock-Technology/uniiotx-querier/interactors"
	"github.com/Bedrock-Technology/uniiotx-querier/logger"
	"github.com/Bedrock-Technology/uniiotx-querier/poller"
	"github.com/Bedrock-Technology/uniiotx-querier/servers"
	"github.com/Bedrock-Technology/uniiotx-querier/storer"

	"github.com/dgraph-io/ristretto"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "uniIOTX contract querier",
	Long:  `An app provides a REST API for querying uniIOTX project's contract status`,
	Run: func(cmd *cobra.Command, args []string) {

		// -------------------------------------------------------------------------------------------------------------
		// Init Program
		// -------------------------------------------------------------------------------------------------------------

		// Create a global context
		ctx, cancel := context.WithCancel(context.Background())

		// Crete logger
		logger.DevMode = config.C.DevMode
		logger.ConsoleEncoder = config.C.ConsoleEncoder
		logger.Stacktrace = config.C.Stacktrace
		logger.RollingFilename = config.C.LogFileName
		myLogger := logger.New(logger.NewZapLogger(logger.DefaultCallerSkip))

		// Print log for check
		myLogger.Info("configs parsed", "configs", config.C)

		// Create storer
		myStorer := storer.NewStorer(config.C.SqliteDSN, myLogger)

		// Create cacher
		c, _ := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,              // number of keys to track frequency of (10M).
			MaxCost:     1 << 30,          // maximum cost of cache (1GB).
			BufferItems: 64,               // number of keys per Get buffer.
			Metrics:     config.C.DevMode, // determines whether cache statistics are kept during the cache's lifetime.
		})

		// Create contract callers
		ethcli, err := ethclient.Dial(config.C.ChainHost)
		if err != nil {
			myLogger.Fatal("failed to create eth client", err)
		}
		iotxStakigCaller, _ := bindings.NewIOTXStakingCaller(ethcommon.HexToAddress(config.C.IOTXStaking), ethcli)
		systemStakigCaller, _ := bindings.NewSystemStakingCaller(ethcommon.HexToAddress(config.C.SystemStaking), ethcli)

		// Create data server
		interactorFactory := &interactors.InteractorFactory{
			Cacher: c,
			Storer: myStorer,
		}
		dataServer := &servers.DataServer{
			Logger: myLogger,
			Addr:   config.C.DataServerAddr,
			If:     interactorFactory,
		}

		// Create poller
		myPoller := poller.Poller{
			Logger:              myLogger,
			SystemStakingCaller: systemStakigCaller,
			IOTXStakingCaller:   iotxStakigCaller,
			Cacher:              c,
			Storer:              myStorer,
		}

		// Create metric server
		metricServer := servers.MetricServer{
			Logger: myLogger,
			Addr:   config.C.MetricServerAddr,
		}

		// -------------------------------------------------------------------------------------------------------------
		// Start Program
		// -------------------------------------------------------------------------------------------------------------

		// Please pay attention: the following goroutines will start in a random order.
		myLogger.Info("starting program", "pid", os.Getppid())

		wg := &sync.WaitGroup{}

		// Listen to the interrupt & termination signals as well as the context cancellation.
		wg.Add(1)
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			myLogger.Info("listening to signals")
		listen:
			for {
				select {
				case s := <-c:
					myLogger.Info("got signal", "signal", s)
					cancel()
					break listen
				case <-ctx.Done():
					myLogger.Info("contex is canceled")
					break listen
				}
			}
			wg.Done()
			myLogger.Info("stopped listening to signals")
		}()

		// Start poller
		wg.Add(1)
		go func() {
			myPoller.Start()
			wg.Done()
		}()

		// Start data server
		wg.Add(1)
		go func() {
			dataServer.Start()
			wg.Done()
		}()

		// Start metric server
		wg.Add(1)
		go func() {
			metricServer.Start()
			wg.Done()
		}()

		// -------------------------------------------------------------------------------------------------------------
		// Exit Program
		// -------------------------------------------------------------------------------------------------------------

		// Shut down services if shutdown is triggerred
		<-ctx.Done()

		ethcli.Close()
		myPoller.Close()
		dataServer.Close()
		metricServer.Close()
		myStorer.Close()
		c.Close()

		wg.Wait()

		if config.C.LogFileName != "" {
			myLogger.Sync()
		}

		myLogger.Info("program exited")
	},
}
