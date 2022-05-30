package functions

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yurchenkosv/metric-service/internal/types"
)

//type metricConstraint interface {
//	*types.Counter | *types.Gauge
//}

//func appendMetric[T metricConstraint](name string, value T, mType string, metrics *types.Metrics) {
//	metric := types.Metric{
//		ID:    name,
//		MType: mType,
//	}
//	switch mType {
//	case "gauge":
//		metric.Value = types.Gauge(value)
//	case "counter":
//		append(metrics.Metric, types.Metric{
//			ID:    name,
//			MType: mType,
//			Delta: types.Counter(value),
//		})
//
//}

var mutex sync.Mutex

func appendGaugeMetric(name string, value float64, metrics *types.Metrics, cfg *types.AgentConfig) {
	gauge := &value
	hash := ""
	if cfg.Key != "" {
		msg := fmt.Sprintf("%s:gauge:%f", name, value)
		hash = CreateSignedHash(msg, []byte(cfg.Key))
	}
	metrics.Metric = append(metrics.Metric, types.Metric{
		ID:    name,
		MType: "gauge",
		Value: gauge,
		Hash:  hash,
	})
}

func appendCounterMetric(name string, value int64, metrics *types.Metrics, cfg *types.AgentConfig) {
	counter := &value
	hash := ""
	if cfg.Key != "" {
		msg := fmt.Sprintf("%s:counter:%d", name, value)
		hash = CreateSignedHash(msg, []byte(cfg.Key))
	}
	metrics.Metric = append(metrics.Metric, types.Metric{
		ID:    name,
		MType: "counter",
		Delta: counter,
		Hash:  hash,
	})
}

func CollectMemMetrics(poolCount int, cfg *types.AgentConfig) types.Metrics {
	var rtm runtime.MemStats
	var memoryMetrics types.Metrics
	runtime.ReadMemStats(&rtm)
	appendGaugeMetric("Alloc", float64(rtm.Alloc), &memoryMetrics, cfg)
	appendGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys), &memoryMetrics, cfg)
	appendGaugeMetric("Frees", float64(rtm.Frees), &memoryMetrics, cfg)
	appendGaugeMetric("GCCPUFraction", float64(rtm.GCCPUFraction), &memoryMetrics, cfg)
	appendGaugeMetric("GCSys", float64(rtm.GCSys), &memoryMetrics, cfg)
	appendGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc), &memoryMetrics, cfg)
	appendGaugeMetric("HeapIdle", float64(rtm.HeapIdle), &memoryMetrics, cfg)
	appendGaugeMetric("HeapInuse", float64(rtm.HeapInuse), &memoryMetrics, cfg)
	appendGaugeMetric("HeapObjects", float64(rtm.HeapObjects), &memoryMetrics, cfg)
	appendGaugeMetric("HeapReleased", float64(rtm.HeapReleased), &memoryMetrics, cfg)
	appendGaugeMetric("HeapSys", float64(rtm.HeapSys), &memoryMetrics, cfg)
	appendGaugeMetric("LastGC", float64(rtm.LastGC), &memoryMetrics, cfg)
	appendGaugeMetric("Lookups", float64(rtm.Lookups), &memoryMetrics, cfg)
	appendGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse), &memoryMetrics, cfg)
	appendGaugeMetric("MCacheSys", float64(rtm.MCacheSys), &memoryMetrics, cfg)
	appendGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse), &memoryMetrics, cfg)
	appendGaugeMetric("MSpanSys", float64(rtm.MSpanSys), &memoryMetrics, cfg)
	appendGaugeMetric("Mallocs", float64(rtm.Mallocs), &memoryMetrics, cfg)
	appendGaugeMetric("NextGC", float64(rtm.NextGC), &memoryMetrics, cfg)
	appendGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC), &memoryMetrics, cfg)
	appendGaugeMetric("NumGC", float64(rtm.NumGC), &memoryMetrics, cfg)
	appendGaugeMetric("OtherSys", float64(rtm.OtherSys), &memoryMetrics, cfg)
	appendGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs), &memoryMetrics, cfg)
	appendGaugeMetric("StackInuse", float64(rtm.StackInuse), &memoryMetrics, cfg)
	appendGaugeMetric("StackSys", float64(rtm.StackSys), &memoryMetrics, cfg)
	appendGaugeMetric("Sys", float64(rtm.Sys), &memoryMetrics, cfg)
	appendGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc), &memoryMetrics, cfg)
	appendGaugeMetric("RandomValue", rand.Float64(), &memoryMetrics, cfg)
	appendCounterMetric("PollCount", int64(poolCount), &memoryMetrics, cfg)
	return memoryMetrics
}

func PushMemMetrics(m types.Metrics, cfg *types.AgentConfig) {
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetBaseURL("http://" + cfg.Address)
	go func() {
		if len(m.Metric) > 0 {
			_, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(m.Metric).
				Post("/updates")
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}

func FlushMetricsToDisk(cfg *types.ServerConfig, m storage.Repository) {
	if cfg.StoreFile == "" {
		return
	}

	fileLocation := cfg.StoreFile
	fileBits := os.O_WRONLY | os.O_CREATE | os.O_TRUNC

	mutex.Lock()
	file, err := os.OpenFile(fileLocation, fileBits, 0600)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(m.AsMetrics())
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()
	mutex.Unlock()
}

func ReadMetricsFromDisk(cnf *types.ServerConfig, repository *storage.Repository) storage.Repository {
	repo := *repository
	fileLocation := cnf.StoreFile

	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Println(err)
		os.Create(fileLocation)
		return repo
	}
	metrics := types.Metrics{}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Println(err)
		return repo
	}

	for i := range metrics.Metric {
		metricName := metrics.Metric[i].ID
		if metrics.Metric[i].MType == "counter" {
			metricValue := metrics.Metric[i].Delta
			repo.AddCounter(metricName, types.Counter(*metricValue))
		}
		if metrics.Metric[i].MType == "gauge" {
			metricValue := metrics.Metric[i].Value
			repo.AddGauge(metricName, types.Gauge(*metricValue))
		}
	}
	return repo
}

func Cleanup(mainLoop *time.Ticker, pushLoop *time.Ticker, mainLoopStop chan bool) {
	mainLoop.Stop()
	pushLoop.Stop()
	mainLoopStop <- true
	println("Program exit")
}

func CreateSignedHash(msg string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(msg))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
