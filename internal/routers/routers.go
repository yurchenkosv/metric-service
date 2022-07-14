package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yurchenkosv/metric-service/internal/handlers"
	"github.com/yurchenkosv/metric-service/internal/middlewares"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"github.com/yurchenkosv/metric-service/internal/types"
)

func NewRouter(cfg *types.ServerConfig, store *storage.Repository) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.AppendConfigToContext(cfg))
	router.Use(middlewares.AddStorage(store))
	router.Use(middlewares.GzipCompress)
	router.Use(middlewares.GzipDecompress)

	router.With(middlewares.SaveMetricToFile).Route("/update", func(r chi.Router) {
		r.Post("/", handlers.HandleUpdateMetricJSON)
		r.Post("/{metricType}/{metricName}/{metricValue}", handlers.HandleUpdateMetric)
	})
	router.Route("/", func(r chi.Router) {
		r.Get("/", handlers.HandleGetAllMetrics)
	})
	router.Route("/value", func(r chi.Router) {
		r.Post("/", handlers.HandleGetMetricJSON)
		r.Get("/{metricType}/{metricName}", handlers.HandleGetMetric)
	})
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", handlers.HealthChecks)
	})
	router.Route("/updates", func(r chi.Router) {
		r.Post("/", handlers.HandleUpdatesJSON)
	})
	return router
}
