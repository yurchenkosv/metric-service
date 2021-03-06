package storage

import (
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/types"
)

type mapStorage struct {
	GaugeMetric   map[string]types.Gauge
	CounterMetric map[string]types.Counter
}

func NewMapStorage() Repository {
	return &mapStorage{
		GaugeMetric:   make(map[string]types.Gauge),
		CounterMetric: make(map[string]types.Counter),
	}
}

func (m *mapStorage) AddCounter(name string, val types.Counter) {
	if len(m.CounterMetric) == 0 {
		m.CounterMetric = make(map[string]types.Counter)
	}
	m.CounterMetric[name] += val
}

func (m *mapStorage) AddGauge(name string, val types.Gauge) {
	if len(m.GaugeMetric) == 0 {
		m.GaugeMetric = make(map[string]types.Gauge)
	}
	m.GaugeMetric[name] = val
}

func (m *mapStorage) GetMetricByKey(key string) (string, error) {
	if val, ok := m.CounterMetric[key]; ok {
		return fmt.Sprintf("%v", val), nil
	}
	if val, ok := m.GaugeMetric[key]; ok {
		return fmt.Sprintf("%.3f", val), nil
	}
	return "", ErrNotFound
}

func (m *mapStorage) GetCounterByKey(key string) (types.Counter, error) {
	if val, ok := m.CounterMetric[key]; ok {
		return val, nil
	}
	return 0, ErrNotFound
}

func (m *mapStorage) GetGaugeByKey(key string) (types.Gauge, error) {
	if val, ok := m.GaugeMetric[key]; ok {
		return val, nil
	}
	return 0, ErrNotFound
}

func (m *mapStorage) GetAllMetrics() string {
	var metrics string
	for k, v := range m.CounterMetric {
		metrics += fmt.Sprintf("key = %s value = %v\n", k, v)
	}
	for k, v := range m.GaugeMetric {
		metrics += fmt.Sprintf("key = %s value = %v\n", k, v)
	}
	return metrics
}

func (m *mapStorage) AsMetrics() types.Metrics {
	var metrics types.Metrics
	for k, v := range m.CounterMetric {
		counter := int64(v)
		metrics.Metric = append(metrics.Metric, types.Metric{
			ID:    k,
			MType: "counter",
			Delta: &counter,
		})
	}
	for k, v := range m.GaugeMetric {
		gauge := float64(v)
		metrics.Metric = append(metrics.Metric, types.Metric{
			ID:    k,
			MType: "gauge",
			Value: &gauge,
		})
	}
	return metrics
}

func (m *mapStorage) InsertMetrics(metrics []types.Metric) {
	for i := range metrics {
		if metrics[i].MType == "counter" {
			counter := *metrics[i].Delta
			m.AddCounter(metrics[i].ID, types.Counter(counter))
		}
		if metrics[i].MType == "gauge" {
			gauge := *metrics[i].Value
			m.AddGauge(metrics[i].ID, types.Gauge(gauge))
		}
	}
}
