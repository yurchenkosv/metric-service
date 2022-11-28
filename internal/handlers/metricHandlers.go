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

// MetricHandler struct that we need for passing service.ServerMetricService into.
type MetricHandler struct {
	metricService *service.ServerMetricService
}

// NewMetricHandler sets service.ServerMetricService and returns pointer to MetricHandler
func NewMetricHandler(metricService *service.ServerMetricService) *MetricHandler {
	return &MetricHandler{metricService: metricService}
}

// validateMetric simple function that return true if metric name is valid and false if not.
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

// HandleUpdateMetricJSON handler for single metric in JSON format.
// It unmarshall metric and pass it to service.ServerMetricService to save.
func (h MetricHandler) HandleUpdateMetricJSON(writer http.ResponseWriter, request *http.Request) {
	var metric model.Metric

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Error("cannot read request", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("receive message with content: '%s'", string(body))
	err = json.Unmarshal(body, &metric)
	if err != nil {
		log.Error("cannot unmarshall", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.metricService.AddMetric(request.Context(), metric)
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

// HandleUpdatesJSON handler for batch of metrics in JSON format.
// It unmarshalls metrics then pass them to service.ServerMetricService to save.
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
	log.Debugf("receive message with content: '%s'", string(data))

	err = json.Unmarshal(data, &metrics)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}

	err = h.metricService.AddMetricBatch(request.Context(), model.Metrics{Metric: metrics})
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

// HandleUpdateMetric handler for single metric, that transmitted over url.
// It parses parameters /metricType/metricName/metricValue ,
// creates Metric object and passes it to  service.ServerMetricService to save
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
	err := h.metricService.AddMetric(request.Context(), metric)
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

// HandleGetMetric handler to get metric from url params.
// It parses url params /metricType/metricName and returns metric value by name in text representation.
func (h MetricHandler) HandleGetMetric(writer http.ResponseWriter, request *http.Request) {
	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")

	writer.Header().Add("Content-Type", "text/plain")

	if !validateMetric(metricType) {
		writer.WriteHeader(http.StatusNotImplemented)
	}

	metric, err := h.metricService.GetMetricByKey(request.Context(), metricName)
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

// HandleGetAllMetrics handler for print all metrics.
// It queries service.ServerMetricService and prints all metrics in text view
func (h MetricHandler) HandleGetAllMetrics(writer http.ResponseWriter, request *http.Request) {
	metrics, err := h.metricService.GetAllMetrics(request.Context())
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	writer.Write([]byte(metrics.String()))
}

// HandleGetMetricJSON handler for get metric based on JSON query.
// It unmarshalls input JSON, then queries service.ServerMetricService for metric by key.
// Then add fields to original metric and marshalls it to JSON
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

	foundMetric, err := h.metricService.GetMetricByKey(request.Context(), metric.ID)
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
