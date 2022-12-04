package api

import (
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/model"
)

func MetricToApiMetric(metric model.Metric) (*Metric, error) {
	var (
		mtype MetricType
	)
	switch metric.MType {
	case "gauge":
		mtype = MetricType_delta
	case "counter":
		mtype = MetricType_counter
	default:
		return nil, fmt.Errorf("unsupported metric type %s", metric.MType)
	}
	if metric.Value == nil {
		return &Metric{
			Id:    metric.ID,
			Mtype: mtype,
			Delta: int64(*metric.Delta),
			Hash:  metric.Hash,
		}, nil
	} else {
		return &Metric{
			Id:    metric.ID,
			Mtype: mtype,
			Value: float32(*metric.Value),
			Hash:  metric.Hash,
		}, nil
	}
}

func ApiMetricToMetric(metric *Metric) (model.Metric, error) {
	return model.Metric{
		ID:    metric.Id,
		MType: metric.Mtype.String(),
		Delta: model.NewCounter(metric.Delta),
		Value: model.NewGauge(float64(metric.Value)),
		Hash:  metric.Hash,
	}, nil
}
