package storage

import (
	"errors"
	"github.com/yurchenkosv/metric-service/internal/types"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	AddCounter(string, types.Counter)
	AddGauge(string, types.Gauge)
	GetMetricByKey(string) (string, error)
	GetCounterByKey(string) (types.Counter, error)
	GetGaugeByKey(string) (types.Gauge, error)
	GetAllMetrics() string
	AsMetrics() types.Metrics
	InsertMetrics([]types.Metric)
}
