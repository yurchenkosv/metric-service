package repository

import (
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
)

// mapStorage repository realization to store metrics in memory.
type mapStorage struct {
	GaugeMetric   map[string]model.Gauge   // GaugeMetric is map for model.Gauge
	CounterMetric map[string]model.Counter // CounterMetric is map for model.Counter
}

// NewMapRepo initializes maps for store metrics and returns pointer to mapStorage.
func NewMapRepo() *mapStorage {
	return &mapStorage{
		GaugeMetric:   make(map[string]model.Gauge),
		CounterMetric: make(map[string]model.Counter),
	}
}

// Migrate do nothing
func (m mapStorage) Migrate(path string) {
}

// SaveCounter just put counter in map
func (m *mapStorage) SaveCounter(name string, val model.Counter) error {
	if len(m.CounterMetric) == 0 {
		m.CounterMetric = make(map[string]model.Counter)
	}
	m.CounterMetric[name] += val
	return nil
}

// SaveGauge just put gauge in map
func (m *mapStorage) SaveGauge(name string, val model.Gauge) error {
	if len(m.GaugeMetric) == 0 {
		m.GaugeMetric = make(map[string]model.Gauge)
	}
	m.GaugeMetric[name] = val
	return nil
}

// GetMetricByKey trying to find key in map and if finds - return pointer to model.Metric with metric by key.
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
		metric.MType = "gauge"
		metric.Value = &val
		return &metric, nil
	}
	return nil, errors.NoSuchMetricError{MetricName: key}
}

// GetAllMetrics iterates over two maps and put all metrics together and returns pointer to model.Metrics.
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

// Ping always returns no error because of maps always available and cannot be unhealthy.
func (m *mapStorage) Ping() error {
	return nil
}

// SaveMetricsBatch iterates over model.Metric slice and save metrics to two maps.
func (m *mapStorage) SaveMetricsBatch(metrics []model.Metric) error {
	for i := range metrics {
		if metrics[i].MType == "counter" {
			counter := *metrics[i].Delta
			err := m.SaveCounter(metrics[i].ID, counter)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if metrics[i].MType == "gauge" {
			gauge := *metrics[i].Value
			err := m.SaveGauge(metrics[i].ID, gauge)
			if err != nil {
				log.Error(err)
				return err
			}
		}
	}
	return nil
}

// Shutdown just do nothing
func (m mapStorage) Shutdown() {

}
