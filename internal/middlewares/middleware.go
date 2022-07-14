package middlewares

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"encoding/json"
	"fmt"
	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"github.com/yurchenkosv/metric-service/internal/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func AppendConfigToContext(config *types.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, types.ContextKey("config"), config)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func AddStorage(store *storage.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, types.ContextKey("storage"), store)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func SaveMetricToFile(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		config := ctx.Value(types.ContextKey("config")).(*types.ServerConfig)
		if config.StoreInterval == 0 {
			store := ctx.Value(types.ContextKey("storage")).(*storage.Repository)
			mapStorage := *store
			functions.FlushMetricsToDisk(config, mapStorage)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

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
		next.ServeHTTP(types.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return http.HandlerFunc(fn)
}

func CheckHash(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		config := ctx.Value(types.ContextKey("config")).(*types.ServerConfig)

		if config.Key != "" {
			var metric types.Metric
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
			if metric.MType == "counter" {
				counter := *metric.Delta
				msg = fmt.Sprintf("%s:counter:%d", metric.ID, counter)
			} else if metric.MType == "gauge" {
				gauge := *metric.Value
				msg = fmt.Sprintf("%s:gauge:%f", metric.ID, gauge)
			}
			hash := functions.CreateSignedHash(msg, []byte(config.Key))
			if !hmac.Equal([]byte(hash), []byte(metric.Hash)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
