package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yurchenkosv/metric-service/internal/handlers"
	"github.com/yurchenkosv/metric-service/internal/middlewares"
	"github.com/yurchenkosv/metric-service/internal/storage"
	"github.com/yurchenkosv/metric-service/internal/types"
)

func NewRouter(cfg *types.Config, store *storage.Repository) chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	router.Use(middlewares.AppendConfigToContext(cfg))
	router.Use(middlewares.AddStorage(store))

	router.Route("/update", func(r chi.Router) {
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
	return router
}
