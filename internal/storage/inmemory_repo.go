package storage

import (
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/types"
)

type Gauge types.Gauge
type Counter types.Counter

type mapStorage struct {
	GaugeMetric   map[string]Gauge
	CounterMetric map[string]Counter
}

func NewMapStorage() Repository {
	return &mapStorage{
		GaugeMetric:   make(map[string]Gauge),
		CounterMetric: make(map[string]Counter),
	}
}

func (m *mapStorage) AddCounter(name string, val Counter) {
	if len(m.CounterMetric) == 0 {
		m.CounterMetric = make(map[string]Counter)
	}
	m.CounterMetric[name] += val
}

func (m *mapStorage) AddGauge(name string, val Gauge) {
	if len(m.GaugeMetric) == 0 {
		m.GaugeMetric = make(map[string]Gauge)
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
