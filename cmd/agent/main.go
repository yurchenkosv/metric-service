package main

import (
	"fmt"
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

type gauge float64
type counter int64

type MemMetrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
	gaugeMetrics  map[string]gauge
}

func collectMemMetrics(poolCount int) MemMetrics {
	var rtm runtime.MemStats
	var memoryMetrics MemMetrics
	metrics := make(map[string]gauge)
	metrics["Alloc"] = gauge(rtm.Alloc)
	metrics["BuckHashSys"] = gauge(rtm.BuckHashSys)
	metrics["Frees"] = gauge(rtm.Frees)
	metrics["GCCPUFraction"] = gauge(rtm.GCCPUFraction)
	metrics["GCSys"] = gauge(rtm.GCSys)
	metrics["HeapAlloc"] = gauge(rtm.HeapAlloc)
	metrics["HeapIdle"] = gauge(rtm.HeapIdle)
	metrics["HeapInuse"] = gauge(rtm.HeapInuse)
	metrics["HeapObjects"] = gauge(rtm.HeapObjects)
	metrics["HeapReleased"] = gauge(rtm.HeapReleased)
	metrics["HeapSys"] = gauge(rtm.HeapSys)
	metrics["LastGC"] = gauge(rtm.LastGC)
	metrics["Lookups"] = gauge(rtm.Lookups)
	metrics["MCacheInuse"] = gauge(rtm.MCacheInuse)
	metrics["MCacheSys"] = gauge(rtm.MCacheSys)
	metrics["MSpanInuse"] = gauge(rtm.MSpanInuse)
	metrics["MSpanSys"] = gauge(rtm.MSpanSys)
	metrics["Mallocs"] = gauge(rtm.Mallocs)
	metrics["NextGC"] = gauge(rtm.NextGC)
	metrics["NumForcedGC"] = gauge(rtm.NumForcedGC)
	metrics["NumGC"] = gauge(rtm.NumGC)
	metrics["OtherSys"] = gauge(rtm.OtherSys)
	metrics["PauseTotalNs"] = gauge(rtm.PauseTotalNs)
	metrics["StackInuse"] = gauge(rtm.StackInuse)
	metrics["StackSys"] = gauge(rtm.StackSys)
	metrics["Sys"] = gauge(rtm.Sys)
	metrics["TotalAlloc"] = gauge(rtm.TotalAlloc)
	metrics["RandomValue"] = gauge(rand.Float64())
	memoryMetrics.PollCount = counter(poolCount)
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
				pollCount += 1
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
