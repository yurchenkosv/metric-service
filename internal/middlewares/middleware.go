package middlewares

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/service"
)

// GzipDecompress is middleware to decompress message body, that was compressed with Gzip algo.
// It executes only if Content-Encoding is gzip.
func GzipDecompress(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer gz.Close()

		body, err := io.ReadAll(gz)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// GzipCompress middleware to compress message body with Gzip algo.
// It executes only if Accept-Encoding is gzip.
func GzipCompress(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(model.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return http.HandlerFunc(fn)
}

// CheckHash middleware to get hash from JSON message and check, that hash is valid.
// It unmarshalls metric and calculate signed hash.
// If hash in message and calculated the same - pass message to handler.
func CheckHash(svc *service.ServerMetricService) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var metric model.Metric
			var msg string
			data, err := io.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewReader(data))
			if err != nil {
				log.Fatal(err)
				return
			}
			err = json.Unmarshal(data, &metric)
			if err != nil {
				log.Fatal(err)
				return
			}
			switch metric.MType {
			case "counter":
				msg = fmt.Sprintf("%s:counter:%s", metric.ID, metric.Delta.String())
			case "gauge":
				msg = fmt.Sprintf("%s:gauge:%s", metric.ID, metric.Value.String())
			}
			hash, err := svc.CreateSignedHash(msg)
			if err != nil {
				log.Error(err)
				return
			}
			if !hmac.Equal([]byte(hash), []byte(metric.Hash)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
