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

	router.Group(func(gr chi.Router) {
		gr.Use(middleware.AllowContentType("text/plain", "application/json"))
		gr.Use(middlewares.GzipCompress)
		gr.Use(middlewares.GzipDecompress)
		gr.Route("/update", func(r chi.Router) {
			r.Post("/", metricHandler.HandleUpdateMetricJSON)
			r.Post("/{metricType}/{metricName}/{metricValue}", metricHandler.HandleUpdateMetric)
		})
		gr.Route("/", func(r chi.Router) {
			r.Get("/", metricHandler.HandleGetAllMetrics)
		})
		gr.Route("/value", func(r chi.Router) {
			r.Post("/", metricHandler.HandleGetMetricJSON)
			r.Get("/{metricType}/{metricName}", metricHandler.HandleGetMetric)
		})
		gr.Route("/ping", func(r chi.Router) {
			r.Get("/", healthCheckHandler.HandleHealthChecks)
		})
		gr.Route("/updates", func(r chi.Router) {
			r.Post("/", metricHandler.HandleUpdatesJSON)
		})
	})
	router.Mount("/debug", middleware.Profiler())
	return router
}
