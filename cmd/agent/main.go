package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	cfg = types.AgentConfig{}
)

func init() {
}

func main() {
	err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}

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
				functions.CollectMemMetrics(pollCount, &cfg)
			case <-pushLoop.C:
				memMetrics <- functions.CollectMemMetrics(pollCount, &cfg)
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
