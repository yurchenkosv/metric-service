package service

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/clients"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
)

// AgentMetricService serves all work with metrics on agent (metric collector) side
type AgentMetricService struct {
	config  *config.AgentConfig
	client  clients.MetricsClient
	metrics *model.Metrics
	mutex   sync.Mutex
}

// NewAgentMetricService creates new AgentMetricService with filled fields and returns pointer on this object
func NewAgentMetricService(cfg *config.AgentConfig, client clients.MetricsClient) *AgentMetricService {
	return &AgentMetricService{
		config:  cfg,
		client:  client,
		metrics: &model.Metrics{Metric: []model.Metric{}},
	}
}

// Push is function to send metrics to server via network
func (s *AgentMetricService) Push() error {
	log.Debug("starting to send metrics to server")
	if len(s.metrics.Metric) == 0 {
		log.Warn("looks like metrics are empty, nothing to send:, ", s.metrics.Metric)
		return nil
	}
	s.client.PushMetrics(s.metrics.Metric)
	s.metrics = &model.Metrics{Metric: []model.Metric{}}
	log.Debug("finish sending metrics to server")
	return nil
}

// CreateSignedHash creates hash signed by key in config. Then this hash adds to metric.
func (s *AgentMetricService) CreateSignedHash(msg string) (string, error) {
	if s.config.HashKey == "" {
		return "", &errors.NoEncryptionKeyFoundError{}
	}
	hash := signHash(s.config.HashKey, msg)
	return hash, nil
}

// CollectMetrics main method for this service. It's collects cpu and RAM metrics
// and stores it in metrics field of AgentMetricService
func (s *AgentMetricService) CollectMetrics(pollCount int) {
	var wg sync.WaitGroup

	wg.Add(2)

	log.Debug("starting new metric poll: ", pollCount)

	go func() {
		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)

		s.appendGaugeMetric("Alloc", float64(rtm.Alloc))
		s.appendGaugeMetric("BuckHashSys", float64(rtm.BuckHashSys))
		s.appendGaugeMetric("Frees", float64(rtm.Frees))
		s.appendGaugeMetric("GCCPUFraction", rtm.GCCPUFraction)
		s.appendGaugeMetric("GCSys", float64(rtm.GCSys))
		s.appendGaugeMetric("HeapAlloc", float64(rtm.HeapAlloc))
		s.appendGaugeMetric("HeapIdle", float64(rtm.HeapIdle))
		s.appendGaugeMetric("HeapInuse", float64(rtm.HeapInuse))
		s.appendGaugeMetric("HeapObjects", float64(rtm.HeapObjects))
		s.appendGaugeMetric("HeapReleased", float64(rtm.HeapReleased))
		s.appendGaugeMetric("HeapSys", float64(rtm.HeapSys))
		s.appendGaugeMetric("LastGC", float64(rtm.LastGC))
		s.appendGaugeMetric("Lookups", float64(rtm.Lookups))
		s.appendGaugeMetric("MCacheInuse", float64(rtm.MCacheInuse))
		s.appendGaugeMetric("MCacheSys", float64(rtm.MCacheSys))
		s.appendGaugeMetric("MSpanInuse", float64(rtm.MSpanInuse))
		s.appendGaugeMetric("MSpanSys", float64(rtm.MSpanSys))
		s.appendGaugeMetric("Mallocs", float64(rtm.Mallocs))
		s.appendGaugeMetric("NextGC", float64(rtm.NextGC))
		s.appendGaugeMetric("NumForcedGC", float64(rtm.NumForcedGC))
		s.appendGaugeMetric("NumGC", float64(rtm.NumGC))
		s.appendGaugeMetric("OtherSys", float64(rtm.OtherSys))
		s.appendGaugeMetric("PauseTotalNs", float64(rtm.PauseTotalNs))
		s.appendGaugeMetric("StackInuse", float64(rtm.StackInuse))
		s.appendGaugeMetric("StackSys", float64(rtm.StackSys))
		s.appendGaugeMetric("Sys", float64(rtm.Sys))
		s.appendGaugeMetric("TotalAlloc", float64(rtm.TotalAlloc))
		s.appendGaugeMetric("RandomValue", rand.Float64())
		s.appendCounterMetric("PollCount", int64(pollCount))
		wg.Done()
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
			s.appendGaugeMetric(metricName, cpuUtil[util])
		}
		s.appendGaugeMetric("TotalMemory", float64(memMetrics.Total))
		s.appendGaugeMetric("FreeMemory", float64(memMetrics.Free))
		wg.Done()
	}()

	wg.Wait()

	log.Debugf("poll %d sucessfuly finished", pollCount)
}

func (s *AgentMetricService) appendGaugeMetric(name string, value float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	gauge := model.NewGauge(value)
	msg := fmt.Sprintf("%s:gauge:%s", name, gauge.String())

	hash, err := s.CreateSignedHash(msg)
	if err != nil {
		log.Info(err)
	}
	s.metrics.Metric = append(s.metrics.Metric, model.Metric{
		ID:    name,
		MType: "gauge",
		Value: model.NewGauge(value),
		Hash:  hash,
	})
}

func (s *AgentMetricService) appendCounterMetric(name string, value int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	counter := model.NewCounter(value)
	msg := fmt.Sprintf("%s:counter:%d", name, value)

	hash, err := s.CreateSignedHash(msg)
	if err != nil {
		log.Info(err)
	}
	s.metrics.Metric = append(s.metrics.Metric, model.Metric{
		ID:    name,
		MType: "counter",
		Delta: counter,
		Hash:  hash,
	})
}
