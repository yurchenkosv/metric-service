package storage

import (
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Repository interface {
	AddCounter(string, Counter)
	AddGauge(string, Gauge)
	GetMetricByKey(string) (string, error)
	GetAllMetrics() string
}
