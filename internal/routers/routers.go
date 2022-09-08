package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/handlers"
	"github.com/yurchenkosv/metric-service/internal/middlewares"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"github.com/yurchenkosv/metric-service/internal/service"
)

func NewRouter(cfg *config.ServerConfig, store repository.Repository) chi.Router {
	var (
		metricService      = service.NewServerMetricService(cfg, store)
		healthCheckService = service.NewHealthCheckService(cfg, store)

		metricHandler      = handlers.NewMetricHandler(metricService)
		healthCheckHandler = handlers.NewHealthCheckHandler(healthCheckService)
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.GzipCompress)
	router.Use(middlewares.GzipDecompress)

	router.Route("/update", func(r chi.Router) {
		r.With(middlewares.CheckHash(metricService)).Post("/", metricHandler.HandleUpdateMetricJSON)
		r.Post("/{metricType}/{metricName}/{metricValue}", metricHandler.HandleUpdateMetric)
	})
	router.Route("/", func(r chi.Router) {
		r.Get("/", metricHandler.HandleGetAllMetrics)
	})
	router.Route("/value", func(r chi.Router) {
		r.Post("/", metricHandler.HandleGetMetricJSON)
		r.Get("/{metricType}/{metricName}", metricHandler.HandleGetMetric)
	})
	router.Route("/ping", func(r chi.Router) {
		r.Get("/", healthCheckHandler.HandleHealthChecks)
	})
	router.Route("/updates", func(r chi.Router) {
		r.Post("/", metricHandler.HandleUpdatesJSON)
	})
	return router
}
