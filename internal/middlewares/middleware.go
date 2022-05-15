package middlewares

import (
	"context"
	"github.com/yurchenkosv/metric-service/internal/functions"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"github.com/yurchenkosv/metric-service/internal/types"
	"net/http"
)

func AppendConfigToContext(config *types.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "config", config)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func AddStorage(store *storage.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "storage", store)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func SaveMetricToFile(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		config := ctx.Value("config").(*types.Config)
		if config.StoreInterval == 0 {
			store := ctx.Value("storage").(*storage.Repository)
			mapStorage := *store
			functions.FlushMetricsToDisk(config, mapStorage)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
