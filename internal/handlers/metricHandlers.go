package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/service"
)

type MetricHandler struct {
	metricService *service.ServerMetricService
}

func NewMetricHandler(metricService *service.ServerMetricService) *MetricHandler {
	return &MetricHandler{metricService: metricService}
}

func validateMetric(metricType string) bool {
	switch metricType {
	case "counter":
		return true
	case "gauge":
		return true
	default:
		return false
	}
}

func (h MetricHandler) HandleUpdateMetricJSON(writer http.ResponseWriter, request *http.Request) {
	var metric model.Metric

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error("cannot read request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Error("cannot unmarshall", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.metricService.AddMetric(metric)
	if err != nil {
		switch e := err.(type) {
		case *errors.NoSuchMetricError:
			writer.WriteHeader(http.StatusNotImplemented)
		default:
			log.Error("unknown error when add metric ", e)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h MetricHandler) HandleUpdatesJSON(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	var metrics []model.Metric

	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	err = h.metricService.AddMetricBatch(model.Metrics{Metric: metrics})
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (h MetricHandler) HandleUpdateMetric(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")
	metricValue := chi.URLParam(request, "metricValue")

	writer.Header().Add("Content-Type", "text/plain")

	metric := model.Metric{
		ID:    metricName,
		MType: metricType,
	}
	switch metricType {
	case "counter":
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			log.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
		}
		metric.Delta = model.NewCounter(val)
	case "gauge":
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			log.Error(err)
			writer.WriteHeader(http.StatusBadRequest)
		}
		metric.Value = model.NewGauge(val)
	default:
		writer.WriteHeader(http.StatusNotImplemented)
	}
	err := h.metricService.AddMetric(metric)
	if err != nil {
		switch err.(type) {
		case *errors.NoSuchMetricError:
			log.Error(err)
			writer.WriteHeader(http.StatusNotImplemented)
		default:
			log.Error(err)
			writer.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (h MetricHandler) HandleGetMetric(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")

	writer.Header().Add("Content-Type", "text/plain")

	if !validateMetric(metricType) {
		writer.WriteHeader(http.StatusNotImplemented)
	}

	metric, err := h.metricService.GetMetricByKey(metricName)
	if err != nil {
		switch err.(type) {
		case *errors.MetricNotFoundError:
			writer.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if metric.Delta != nil {
		writer.Write([]byte(metric.Delta.String()))
	} else if metric.Value != nil {
		writer.Write([]byte(metric.Value.String()))
	}
}

func (h MetricHandler) HandleGetAllMetrics(writer http.ResponseWriter, request *http.Request) {
	metrics, err := h.metricService.GetAllMetrics()
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	writer.Write([]byte(metrics.String()))
}

func (h MetricHandler) HandleGetMetricJSON(writer http.ResponseWriter, request *http.Request) {
	var metric model.Metric
	var msg string

	if request.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &metric)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	foundMetric, err := h.metricService.GetMetricByKey(metric.ID)
	if err != nil {
		switch err.(type) {
		case *errors.MetricNotFoundError:
			writer.WriteHeader(http.StatusNotFound)
			return
		default:
			log.Error(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	switch metric.MType {
	case "counter":
		msg = fmt.Sprintf("%s:counter:%d", foundMetric.ID, *foundMetric.Delta)
	case "gauge":
		msg = fmt.Sprintf("%s:gauge:%f", foundMetric.ID, *foundMetric.Value)
	}

	foundMetric.Hash, err = h.metricService.CreateSignedHash(msg)
	if err != nil {
		log.Info(err)
	}

	data, err = json.Marshal(foundMetric)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(data)
}
