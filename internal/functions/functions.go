package functions

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/yurchenkosv/metric-service/internal/types"
	"log"
	"math/rand"
	"runtime"
	"time"
)

var (
	server = types.URLServer{}
)

func CollectMemMetrics(poolCount int) types.MemMetrics {
	var rtm runtime.MemStats
	var memoryMetrics types.MemMetrics
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
	memoryMetrics.GaugeMetrics = metrics
	return memoryMetrics
}

func PushMemMetrics(m types.MemMetrics) {
	apiServer := server.
		SetHost("localhost").
		SetPort("8080").
		SetSchema("http").
		Build()

	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetBaseURL(apiServer)

	for metricName, metricValue := range m.GaugeMetrics {
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(
				fmt.Sprintf(
					"/update/gauge/%s/%v",
					metricName,
					metricValue,
				),
			)
		if err != nil {
			log.Panic(err)
		}
	}
	_, err := client.R().
		SetHeader("Content-Type", "text/plain").
		Post(
			fmt.Sprintf(
				"/update/counter/%s/%v",
				"PollCount",
				m.PollCount,
			),
		)
	if err != nil {
		log.Panic(err)
	}
}

func Cleanup(mainLoop *time.Ticker, pushLoop *time.Ticker, mainLoopStop chan bool) {
	mainLoop.Stop()
	pushLoop.Stop()
	mainLoopStop <- true
	println("Program exit")
}
