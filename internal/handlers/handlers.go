package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"github.com/yurchenkosv/metric-service/internal/types"
	"io"
	"net/http"
	"strconv"
	"sync"
)

var mutex sync.Mutex

func checkMetricType(metricType string, w http.ResponseWriter) {
	if metricType != "counter" && metricType != "gauge" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}

func checkForError(err error) {
	if err != nil {
		panic(err)
	}
}

func HandleUpdateMetricJSON(writer http.ResponseWriter, request *http.Request) {
	var metrics types.Metric
	ctx := request.Context()
	store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
	mapStorage := *store

	body, err := io.ReadAll(request.Body)
	checkForError(err)
	err = json.Unmarshal(body, &metrics)
	checkForError(err)

	metricType := metrics.MType
	checkMetricType(metricType, writer)
	mutex.Lock()
	defer mutex.Unlock()
	if metricType == "counter" {
		counter := types.Counter(*metrics.Delta)
		mapStorage.AddCounter(metrics.ID, counter)
	}
	if metricType == "gauge" {
		gauge := types.Gauge(*metrics.Value)
		mapStorage.AddGauge(metrics.ID, gauge)
	}
}

func HandleUpdateMetric(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
	mapStorage := *store

	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")
	metricValue := chi.URLParam(request, "metricValue")

	checkMetricType(metricType, writer)
	mutex.Lock()
	defer mutex.Unlock()
	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		mapStorage.AddCounter(metricName, types.Counter(val))
	}
	if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		mapStorage.AddGauge(metricName, types.Gauge(val))
	}
}

func HandleGetMetric(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
	mapStorage := *store

	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")
	checkMetricType(metricType, writer)

	val, err := mapStorage.GetMetricByKey(metricName)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("no metrics found"))
	}
	writer.Write([]byte(val))
}

func HandleGetAllMetrics(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
	mapStorage := *store

	val := mapStorage.GetAllMetrics()
	writer.Header().Set("Content-Type", "text/html")
	writer.Write([]byte(val))
}

func HandleGetMetricJSON(writer http.ResponseWriter, request *http.Request) {
	var metric types.Metric
	if request.Header.Get("Content-Type") != "application/json" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := request.Context()
	store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
	mapStorage := *store

	data, err := io.ReadAll(request.Body)
	checkForError(err)
	err = json.Unmarshal(data, &metric)
	checkForError(err)

	if metric.MType == "counter" {
		val, err := mapStorage.GetCounterByKey(metric.ID)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		counter := int64(val)
		metric.Delta = &counter
	}
	if metric.MType == "gauge" {
		val, err := mapStorage.GetGaugeByKey(metric.ID)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		gauge := float64(val)
		metric.Value = &gauge
	}

	data, err = json.Marshal(metric)
	checkForError(err)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(data)
}
