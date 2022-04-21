package storage

type Repository interface {
	Save() bool
}

type Gauge float64
type Counter int64

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
