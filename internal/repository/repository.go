package repository

import (
	"github.com/yurchenkosv/metric-service/internal/model"
)

type Repository interface {
	SaveCounter(string, model.Counter) error
	SaveGauge(string, model.Gauge) error
	GetMetricByKey(string) (*model.Metric, error)
	SaveMetricsBatch([]model.Metric) error
	GetAllMetrics() (*model.Metrics, error)
	Ping() error
}
