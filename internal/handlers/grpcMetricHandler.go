package handlers

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/api"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/service"
)

// GRPCMetricHandler struct that we need for passing service.ServerMetricService into.
type GRPCMetricHandler struct {
	metricService *service.ServerMetricService
}

// NewGRPCMetricHandler sets service.ServerMetricService and returns pointer to GRPCMetricHandler
func NewGRPCMetricHandler(metricService *service.ServerMetricService) *GRPCMetricHandler {
	return &GRPCMetricHandler{metricService: metricService}
}

func (h *GRPCMetricHandler) GetMetricByID(ctx context.Context, req *api.MetricRequestByID) (*api.Metric, error) {
	metric, err := h.metricService.GetMetricByKey(ctx, req.GetId())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return api.MetricToAPIMetric(*metric)
}
func (h *GRPCMetricHandler) GetAllMetrics(ctx context.Context, req *api.MetricRequestAll) (*api.Metrics, error) {
	apiMetrics := &api.Metrics{}
	metrics, err := h.metricService.GetAllMetrics(ctx)
	if err != nil {
		return nil, err
	}
	for _, metric := range metrics.Metric {
		apiMetric, err2 := api.MetricToAPIMetric(metric)
		if err2 != nil {
			return nil, err2
		}
		apiMetrics.Metrics = append(apiMetrics.Metrics, apiMetric)
	}
	return apiMetrics, nil
}
func (h *GRPCMetricHandler) SaveMetric(ctx context.Context, req *api.Metric) (*api.MetricResponse, error) {
	metric, err := api.APIMetricToMetric(req)
	if err != nil {
		return &api.MetricResponse{Status: api.MetricStatus_rejected}, err
	}
	err = h.metricService.AddMetric(ctx, metric)
	if err != nil {
		return &api.MetricResponse{Status: api.MetricStatus_rejected}, err
	}
	return &api.MetricResponse{Status: api.MetricStatus_accepted}, nil
}
func (h *GRPCMetricHandler) SaveMetrics(ctx context.Context, req *api.Metrics) (*api.MetricResponse, error) {
	var metrics model.Metrics
	for _, apiMetric := range req.Metrics {
		metric, err := api.APIMetricToMetric(apiMetric)
		if err != nil {
			return &api.MetricResponse{Status: api.MetricStatus_rejected}, err
		}
		metrics.Metric = append(metrics.Metric, metric)
	}
	log.Debug("got metrics ", metrics)
	err := h.metricService.AddMetricBatch(ctx, metrics)
	if err != nil {
		log.Error(err)
		return &api.MetricResponse{Status: api.MetricStatus_rejected}, err
	}
	return &api.MetricResponse{Status: api.MetricStatus_accepted}, nil
}
