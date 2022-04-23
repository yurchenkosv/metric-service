package storage

import "github.com/yurchenkosv/metric-service/internal/types"

type Repository interface {
	Save() bool
}

type Gauge types.Gauge
type Counter types.Counter

type MapStorage struct {
	GaugeMetric   map[string]Gauge
	CounterMetric map[string]Counter
}

func (m *MapStorage) AddCounter(name string, val Counter) {
	if len(m.CounterMetric) == 0 {
		m.CounterMetric = make(map[string]Counter)
	}
	m.CounterMetric[name] += val
}

func (m *MapStorage) AddGauge(name string, val Gauge) {
	if len(m.GaugeMetric) == 0 {
		m.GaugeMetric = make(map[string]Gauge)
	}
	m.GaugeMetric[name] = val
}
