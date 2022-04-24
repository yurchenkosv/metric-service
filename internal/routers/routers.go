package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yurchenkosv/metric-service/internal/handlers"
)

func NewRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/update/{metricType}", func(r chi.Router) {
		r.Post("/{metricName}/{metricValue}", handlers.HandleMetric)

	})
	return router
}