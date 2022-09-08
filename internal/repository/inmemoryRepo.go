package repository

import (
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
)

type mapStorage struct {
	GaugeMetric   map[string]model.Gauge
	CounterMetric map[string]model.Counter
}

func NewMapStorage() *mapStorage {
	return &mapStorage{
		GaugeMetric:   make(map[string]model.Gauge),
		CounterMetric: make(map[string]model.Counter),
	}
}

func (m *mapStorage) SaveCounter(name string, val model.Counter) error {
	if len(m.CounterMetric) == 0 {
		m.CounterMetric = make(map[string]model.Counter)
	}
	m.CounterMetric[name] += val
	return nil
}

func (m *mapStorage) SaveGauge(name string, val model.Gauge) error {
	if len(m.GaugeMetric) == 0 {
		m.GaugeMetric = make(map[string]model.Gauge)
	}
	m.GaugeMetric[name] = val
	return nil
}

func (m *mapStorage) GetMetricByKey(key string) (*model.Metric, error) {
	var metric model.Metric
	if val, ok := m.CounterMetric[key]; ok {
		metric.ID = key
		metric.MType = "counter"
		metric.Delta = &val
		return &metric, nil
	}
	if val, ok := m.GaugeMetric[key]; ok {
		metric.ID = key
		metric.MType = "counter"
		metric.Value = &val
		return &metric, nil
	}
	return nil, errors.NoSuchMetricError{MetricName: key}
}

func (m *mapStorage) GetAllMetrics() (*model.Metrics, error) {
	var metrics model.Metrics
	for k, v := range m.CounterMetric {
		metric := model.Metric{
			ID:    k,
			MType: "counter",
			Delta: &v,
		}
		metrics.Metric = append(metrics.Metric, metric)
	}
	for k, v := range m.GaugeMetric {
		metric := model.Metric{
			ID:    k,
			MType: "gauge",
			Value: &v,
		}
		metrics.Metric = append(metrics.Metric, metric)
	}
	return &metrics, nil
}

func (m *mapStorage) Ping() error {
	return nil
}

//func (m *mapStorage) AsMetrics() types.Metrics {
//	var metrics types.Metrics
//	for k, v := range m.CounterMetric {
//		counter := int64(v)
//		metrics.Metric = append(metrics.Metric, types.Metric{
//			ID:    k,
//			MType: "counter",
//			Delta: &counter,
//		})
//	}
//	for k, v := range m.GaugeMetric {
//		gauge := float64(v)
//		metrics.Metric = append(metrics.Metric, types.Metric{
//			ID:    k,
//			MType: "gauge",
//			Value: &gauge,
//		})
//	}
//	return metrics
//}

//func (m *mapStorage) InsertMetrics(metrics []types.Metric) {
//	for i := range metrics {
//		if metrics[i].MType == "counter" {
//			counter := *metrics[i].Delta
//			m.AddCounter(metrics[i].ID, types.Counter(counter))
//		}
//		if metrics[i].MType == "gauge" {
//			gauge := *metrics[i].Value
//			m.AddGauge(metrics[i].ID, types.Gauge(gauge))
//		}
//	}
//}
