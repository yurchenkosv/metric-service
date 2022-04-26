package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/yurchenkosv/metric-service/internal/storage"
)

var mapStorage = storage.NewMapStorage()

func checkMetricType(metricType string, w http.ResponseWriter) {
	if metricType != "counter" && metricType != "gauge" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
}

func HandleUpdateMetric(writer http.ResponseWriter, request *http.Request) {
	//if request.Header.Get("Content-Type") != "text/plain" {
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")
	metricValue := chi.URLParam(request, "metricValue")

	checkMetricType(metricType, writer)

	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		mapStorage.AddCounter(metricName, storage.Counter(val))
	}
	if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		mapStorage.AddGauge(metricName, storage.Gauge(val))
	}
}

func HandleGetMetric(writer http.ResponseWriter, request *http.Request) {
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
	val := mapStorage.GetAllMetrics()
	writer.Write([]byte(val))
}
