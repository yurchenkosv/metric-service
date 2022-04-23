package handlers

import (
	"github.com/yurchenkosv/metric-service/internal/storage"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var mapStorage = &storage.MapStorage{}

func HandleMetric(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqUrl := r.URL.Path
	urlMatch, _ := regexp.MatchString("update/(gauge|counter)/[a-zA-z]+/(\\d+|\\d+\\.)\\d*$", reqUrl)

	if !urlMatch {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metrics := strings.Split(reqUrl, "/")
	if metrics[len(metrics)-3] == "counter" {
		metricName := metrics[len(metrics)-2]
		metricValue, _ := strconv.ParseInt(metrics[len(metrics)-1], 10, 64)
		mapStorage.AddCounter(metricName, storage.Counter(metricValue))
	} else if metrics[len(metrics)-3] == "gauge" {
		metricName := metrics[len(metrics)-2]
		metricValue, _ := strconv.ParseFloat(metrics[len(metrics)-1], 64)
		mapStorage.AddGauge(metricName, storage.Gauge(metricValue))
	}
	w.Header().Add("Content-Type", "text/plain")
}
