package clients

import "github.com/yurchenkosv/metric-service/internal/model"

// MetricsClient interface to simplify testing and create ability to build several clients for different purposes.
type MetricsClient interface {
	// PushMetrics method that should send metrics to metric server
	PushMetrics(metrics []model.Metric)
}
