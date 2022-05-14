package main

import (
	"github.com/caarlos0/env/v6"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg = types.Config{}
)

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mainLoop := time.NewTicker(cfg.PollInterval)
	pushLoop := time.NewTicker(cfg.ReportInterval)
	mainLoopStop := make(chan bool)
	memMetrics := make(chan types.Metrics)
	osSignal := make(chan os.Signal, 3)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		var pollCount int
		for {
			select {
			case <-mainLoopStop:
				return
			case <-mainLoop.C:
				pollCount = 1
				functions.CollectMemMetrics(pollCount)
			case <-pushLoop.C:
				memMetrics <- functions.CollectMemMetrics(pollCount)
			}
		}
	}()

	go func() {
		for {
			functions.PushMemMetrics(<-memMetrics, &cfg)
		}
	}()

	<-osSignal
	functions.Cleanup(mainLoop, pushLoop, mainLoopStop)
	os.Exit(0)
}
