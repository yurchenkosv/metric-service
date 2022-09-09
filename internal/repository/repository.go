package repository

import (
	"github.com/yurchenkosv/metric-service/internal/model"
)

// Repository interface to gain ability to test service layer.
// It also provides contract to change data storage level of application.
type Repository interface {
	// SaveCounter for  saving model.Counter in storage.
	SaveCounter(string, model.Counter) error

	// SaveGauge for saving model.Gauge in storage.
	SaveGauge(string, model.Gauge) error

	// GetMetricByKey for getting pointer to model.Metric from string key.
	GetMetricByKey(string) (*model.Metric, error)

	// SaveMetricsBatch for saving slice of model.Metric in repository.
	SaveMetricsBatch([]model.Metric) error

	// GetAllMetrics for getting pointer to model.Metrics with all metrics, stored in repository
	GetAllMetrics() (*model.Metrics, error)

	// Shutdown method for graceful shutdown.
	//When it's called, repository should save metrics, close connections and be ready to application shutdown
	Shutdown()

	// Ping should return error when repository assumed as unhealthy
	Ping() error
}
