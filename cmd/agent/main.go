package main

import (
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/types"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type MemMetrics struct {
	Alloc         types.Gauge
	BuckHashSys   types.Gauge
	Frees         types.Gauge
	GCCPUFraction types.Gauge
	GCSys         types.Gauge
	HeapAlloc     types.Gauge
	HeapIdle      types.Gauge
	HeapInuse     types.Gauge
	HeapObjects   types.Gauge
	HeapReleased  types.Gauge
	HeapSys       types.Gauge
	LastGC        types.Gauge
	Lookups       types.Gauge
	MCacheInuse   types.Gauge
	MCacheSys     types.Gauge
	MSpanInuse    types.Gauge
	MSpanSys      types.Gauge
	Mallocs       types.Gauge
	NextGC        types.Gauge
	NumForcedGC   types.Gauge
	NumGC         types.Gauge
	OtherSys      types.Gauge
	PauseTotalNs  types.Gauge
	StackInuse    types.Gauge
	StackSys      types.Gauge
	Sys           types.Gauge
	TotalAlloc    types.Gauge
	PollCount     types.Counter
	RandomValue   types.Gauge
	gaugeMetrics  map[string]types.Gauge
}

func collectMemMetrics(poolCount int) MemMetrics {
	var rtm runtime.MemStats
	var memoryMetrics MemMetrics
	metrics := make(map[string]types.Gauge)
	metrics["Alloc"] = types.Gauge(rtm.Alloc)
	metrics["BuckHashSys"] = types.Gauge(rtm.BuckHashSys)
	metrics["Frees"] = types.Gauge(rtm.Frees)
	metrics["GCCPUFraction"] = types.Gauge(rtm.GCCPUFraction)
	metrics["GCSys"] = types.Gauge(rtm.GCSys)
	metrics["HeapAlloc"] = types.Gauge(rtm.HeapAlloc)
	metrics["HeapIdle"] = types.Gauge(rtm.HeapIdle)
	metrics["HeapInuse"] = types.Gauge(rtm.HeapInuse)
	metrics["HeapObjects"] = types.Gauge(rtm.HeapObjects)
	metrics["HeapReleased"] = types.Gauge(rtm.HeapReleased)
	metrics["HeapSys"] = types.Gauge(rtm.HeapSys)
	metrics["LastGC"] = types.Gauge(rtm.LastGC)
	metrics["Lookups"] = types.Gauge(rtm.Lookups)
	metrics["MCacheInuse"] = types.Gauge(rtm.MCacheInuse)
	metrics["MCacheSys"] = types.Gauge(rtm.MCacheSys)
	metrics["MSpanInuse"] = types.Gauge(rtm.MSpanInuse)
	metrics["MSpanSys"] = types.Gauge(rtm.MSpanSys)
	metrics["Mallocs"] = types.Gauge(rtm.Mallocs)
	metrics["NextGC"] = types.Gauge(rtm.NextGC)
	metrics["NumForcedGC"] = types.Gauge(rtm.NumForcedGC)
	metrics["NumGC"] = types.Gauge(rtm.NumGC)
	metrics["OtherSys"] = types.Gauge(rtm.OtherSys)
	metrics["PauseTotalNs"] = types.Gauge(rtm.PauseTotalNs)
	metrics["StackInuse"] = types.Gauge(rtm.StackInuse)
	metrics["StackSys"] = types.Gauge(rtm.StackSys)
	metrics["Sys"] = types.Gauge(rtm.Sys)
	metrics["TotalAlloc"] = types.Gauge(rtm.TotalAlloc)
	metrics["RandomValue"] = types.Gauge(rand.Float64())
	memoryMetrics.PollCount = types.Counter(poolCount)
	memoryMetrics.gaugeMetrics = metrics
	return memoryMetrics
}

func pushMemMetrics(m MemMetrics) {
	server := "http://localhost:8080"
	client := &http.Client{}
	for metricName, metricValue := range m.gaugeMetrics {
		resp, err := client.Post(fmt.Sprintf("%s/update/gauge/%s/%v", server, metricName, fmt.Sprintf("%f", metricValue)),
			"text/plain", nil)
		if err != nil {
			log.Panic(err)
		}
		defer resp.Body.Close()
	}
	resp, err := client.Post(fmt.Sprintf("%s/update/counter/%s/%v", server, "PollCount", fmt.Sprintf("%v", m.PollCount)),
		"text/plain", nil)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()
}

func cleanup(mainLoop *time.Ticker, pushLoop *time.Ticker, mainLoopStop chan bool) {
	mainLoop.Stop()
	pushLoop.Stop()
	mainLoopStop <- true
	println("Program exit")
}

func main() {
	mainLoop := time.NewTicker(2 * time.Second)
	pushLoop := time.NewTicker(10 * time.Second)
	mainLoopStop := make(chan bool)
	memMetrics := make(chan MemMetrics)
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
				collectMemMetrics(pollCount)
			case <-pushLoop.C:
				memMetrics <- collectMemMetrics(pollCount)
			}
		}
	}()

	go func() {
		for {
			pushMemMetrics(<-memMetrics)
		}
	}()

	go func() {
		<-osSignal
		cleanup(mainLoop, pushLoop, mainLoopStop)
		os.Exit(0)
	}()
	wg.Wait()
}
