package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"log"
	"net/http"
	"strconv"
)

var mapStorage = &storage.MapStorage{}

func HandleMetric(writer http.ResponseWriter, request *http.Request) {
	//if request.Header.Get("Content-Type") != "text/plain" {
	//	writer.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	metricType := chi.URLParam(request, "metricType")
	metricName := chi.URLParam(request, "metricName")
	metricValue := chi.URLParam(request, "metricValue")
	if metricType != "counter" && metricType != "gauge" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if metricType == "counter" {
		val, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			writer.WriteHeader(400)
			log.Panic(err)
		}
		mapStorage.AddCounter(metricName, storage.Counter(val))
	}
	if metricType == "gauge" {
		val, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			writer.WriteHeader(400)
			log.Panic(err)
		}
		mapStorage.AddGauge(metricName, storage.Gauge(val))
	}
}
