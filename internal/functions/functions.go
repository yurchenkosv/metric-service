package functions

import (
	"encoding/json"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"log"
	"math/rand"
	"os"
	"runtime"
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
func appendGaugeMetric(name string, value float64, metrics *types.Metrics) {
	gauge := &value
	metrics.Metric = append(metrics.Metric, types.Metric{
		ID:    name,
		MType: "gauge",
		Value: gauge,
	})
}

func appendCounterMetric(name string, value int64, metrics *types.Metrics) {
	counter := &value
	metrics.Metric = append(metrics.Metric, types.Metric{
		ID:    name,
		MType: "counter",
		Delta: counter,
	})
}

func CollectMemMetrics(poolCount int) types.Metrics {
	var rtm runtime.MemStats
	var memoryMetrics types.Metrics
	runtime.ReadMemStats(&rtm)
	appendGaugeMetric("Alloc", float64(rtm.Alloc), &memoryMetrics)
	appendGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys), &memoryMetrics)
	appendGaugeMetric("Frees", float64(rtm.Frees), &memoryMetrics)
	appendGaugeMetric("GCCPUFraction", float64(rtm.GCCPUFraction), &memoryMetrics)
	appendGaugeMetric("GCSys", float64(rtm.GCSys), &memoryMetrics)
	appendGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc), &memoryMetrics)
	appendGaugeMetric("HeapIdle", float64(rtm.HeapIdle), &memoryMetrics)
	appendGaugeMetric("HeapInuse", float64(rtm.HeapInuse), &memoryMetrics)
	appendGaugeMetric("HeapObjects", float64(rtm.HeapObjects), &memoryMetrics)
	appendGaugeMetric("HeapReleased", float64(rtm.HeapReleased), &memoryMetrics)
	appendGaugeMetric("HeapSys", float64(rtm.HeapSys), &memoryMetrics)
	appendGaugeMetric("LastGC", float64(rtm.LastGC), &memoryMetrics)
	appendGaugeMetric("Lookups", float64(rtm.Lookups), &memoryMetrics)
	appendGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse), &memoryMetrics)
	appendGaugeMetric("MCacheSys", float64(rtm.MCacheSys), &memoryMetrics)
	appendGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse), &memoryMetrics)
	appendGaugeMetric("MSpanSys", float64(rtm.MSpanSys), &memoryMetrics)
	appendGaugeMetric("Mallocs", float64(rtm.Mallocs), &memoryMetrics)
	appendGaugeMetric("NextGC", float64(rtm.NextGC), &memoryMetrics)
	appendGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC), &memoryMetrics)
	appendGaugeMetric("NumGC", float64(rtm.NumGC), &memoryMetrics)
	appendGaugeMetric("OtherSys", float64(rtm.OtherSys), &memoryMetrics)
	appendGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs), &memoryMetrics)
	appendGaugeMetric("StackInuse", float64(rtm.StackInuse), &memoryMetrics)
	appendGaugeMetric("StackSys", float64(rtm.StackSys), &memoryMetrics)
	appendGaugeMetric("Sys", float64(rtm.Sys), &memoryMetrics)
	appendGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc), &memoryMetrics)
	appendGaugeMetric("RandomValue", rand.Float64(), &memoryMetrics)
	appendCounterMetric("PollCount", int64(poolCount), &memoryMetrics)
	return memoryMetrics
}

func PushMemMetrics(m types.Metrics, cfg *types.Config) {
	client := resty.New()
	client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetBaseURL("http://" + cfg.Address)

	for i := range m.Metric {
		metric := m.Metric[i]
		go func() {
			_, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(metric).
				Post("/update")
			if err != nil {
				log.Panic(err)
			}
		}()
	}
}

func FlushMetricsToDisk(cfg *types.Config, m storage.Repository) {

	fileLocation := cfg.StoreFile
	fileBits := os.O_WRONLY | os.O_CREATE

	file, err := os.OpenFile(fileLocation, fileBits, 0600)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.Marshal(m.AsMetrics())

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

//func ReadMetricsFromDisk(cnf *types.Config) types.Metrics {
//	fileLocation := cnf.StoreFile
//}

func Cleanup(mainLoop *time.Ticker, pushLoop *time.Ticker, mainLoopStop chan bool) {
	mainLoop.Stop()
	pushLoop.Stop()
	mainLoopStop <- true
	println("Program exit")
}
