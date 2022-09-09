package clients

import "github.com/yurchenkosv/metric-service/internal/model"

type MetricsClient interface {
	PushMetrics(metrics model.Metrics)
}
