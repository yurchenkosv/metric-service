package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/types"
)

func main() {
	mainLoop := time.NewTicker(2 * time.Second)
	pushLoop := time.NewTicker(10 * time.Second)
	mainLoopStop := make(chan bool)
	memMetrics := make(chan types.MemMetrics)
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	wg := sync.WaitGroup{}
	wg.Add(1)

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
			functions.PushMemMetrics(<-memMetrics)
		}
	}()

	go func() {
		<-osSignal
		functions.Cleanup(mainLoop, pushLoop, mainLoopStop)
		os.Exit(0)
	}()
	wg.Wait()
}
