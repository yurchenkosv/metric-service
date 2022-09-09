package service

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
	"math/rand"
	"runtime"
	"sync"
)

type AgentMetricService struct {
	config *config.AgentConfig
	mutex  sync.Mutex
}

func NewAgentMetricService(cfg *config.AgentConfig) *AgentMetricService {
	return &AgentMetricService{
		config: cfg,
	}
}

func (s *AgentMetricService) CreateSignedHash(msg string) (string, error) {
	if s.config.HashKey == "" {
		return "", &errors.NoEncryptionKeyFoundError{}
	}
	hash := signHash(s.config.HashKey, msg)
	return hash, nil
}

func (s *AgentMetricService) CollectMetrics(poolCount int) model.Metrics {
	var memoryMetrics model.Metrics
	//var wg sync.WaitGroup

	//wg.Add(2)

	go func() {
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)

		s.appendGaugeMetric("Alloc", float64(rtm.Alloc), &memoryMetrics)
		s.appendGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys), &memoryMetrics)
		s.appendGaugeMetric("Frees", float64(rtm.Frees), &memoryMetrics)
		s.appendGaugeMetric("GCCPUFraction", rtm.GCCPUFraction, &memoryMetrics)
		s.appendGaugeMetric("GCSys", float64(rtm.GCSys), &memoryMetrics)
		s.appendGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc), &memoryMetrics)
		s.appendGaugeMetric("HeapIdle", float64(rtm.HeapIdle), &memoryMetrics)
		s.appendGaugeMetric("HeapInuse", float64(rtm.HeapInuse), &memoryMetrics)
		s.appendGaugeMetric("HeapObjects", float64(rtm.HeapObjects), &memoryMetrics)
		s.appendGaugeMetric("HeapReleased", float64(rtm.HeapReleased), &memoryMetrics)
		s.appendGaugeMetric("HeapSys", float64(rtm.HeapSys), &memoryMetrics)
		s.appendGaugeMetric("LastGC", float64(rtm.LastGC), &memoryMetrics)
		s.appendGaugeMetric("Lookups", float64(rtm.Lookups), &memoryMetrics)
		s.appendGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse), &memoryMetrics)
		s.appendGaugeMetric("MCacheSys", float64(rtm.MCacheSys), &memoryMetrics)
		s.appendGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse), &memoryMetrics)
		s.appendGaugeMetric("MSpanSys", float64(rtm.MSpanSys), &memoryMetrics)
		s.appendGaugeMetric("Mallocs", float64(rtm.Mallocs), &memoryMetrics)
		s.appendGaugeMetric("NextGC", float64(rtm.NextGC), &memoryMetrics)
		s.appendGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC), &memoryMetrics)
		s.appendGaugeMetric("NumGC", float64(rtm.NumGC), &memoryMetrics)
		s.appendGaugeMetric("OtherSys", float64(rtm.OtherSys), &memoryMetrics)
		s.appendGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs), &memoryMetrics)
		s.appendGaugeMetric("StackInuse", float64(rtm.StackInuse), &memoryMetrics)
		s.appendGaugeMetric("StackSys", float64(rtm.StackSys), &memoryMetrics)
		s.appendGaugeMetric("Sys", float64(rtm.Sys), &memoryMetrics)
		s.appendGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc), &memoryMetrics)
		s.appendGaugeMetric("RandomValue", rand.Float64(), &memoryMetrics)
		s.appendCounterMetric("PollCount", int64(poolCount), &memoryMetrics)
		//wg.Done()
	}()

	go func() {
		memMetrics, err := mem.VirtualMemory()
		if err != nil {
			log.Error("could not get mem metrics")
		}

		cpuUtil, err := cpu.Percent(0, true)
		if err != nil {
			log.Error("could not get cpu metrics")
		}

		for util := range cpuUtil {
			metricName := fmt.Sprintf("CPUutilization%d", util+1)
			s.appendGaugeMetric(metricName, cpuUtil[util], &memoryMetrics)
		}
		s.appendGaugeMetric("TotalMemory", float64(memMetrics.Total), &memoryMetrics)
		s.appendGaugeMetric("FreeMemory", float64(memMetrics.Free), &memoryMetrics)
		//wg.Done()
	}()

	//wg.Wait()
	return memoryMetrics
}

func (s *AgentMetricService) appendGaugeMetric(name string, value float64, metrics *model.Metrics) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	gauge := model.NewGauge(value)
	msg := fmt.Sprintf("%s:gauge:%s", name, gauge.String())

	hash, err := s.CreateSignedHash(msg)
	if err != nil {
		log.Info(err)
	}
	metrics.Metric = append(metrics.Metric, model.Metric{
		ID:    name,
		MType: "gauge",
		Value: model.NewGauge(value),
		Hash:  hash,
	})
}

func (s *AgentMetricService) appendCounterMetric(name string, value int64, metrics *model.Metrics) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	counter := model.NewCounter(value)
	msg := fmt.Sprintf("%s:counter:%d", name, value)

	hash, err := s.CreateSignedHash(msg)
	if err != nil {
		log.Info(err)
	}
	metrics.Metric = append(metrics.Metric, model.Metric{
		ID:    name,
		MType: "counter",
		Delta: counter,
		Hash:  hash,
	})
}
