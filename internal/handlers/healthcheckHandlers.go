package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/service"
)

// HealthChecksHandler struct to pass service.HealthCheckService into.
type HealthChecksHandler struct {
	svc *service.HealthCheckService
}

// NewHealthCheckHandler sets service.HealthCheckService and returns pointer to HealthChecksHandler.
func NewHealthCheckHandler(svc *service.HealthCheckService) *HealthChecksHandler {
	return &HealthChecksHandler{
		svc: svc,
	}
}

// HandleHealthChecks handler for checking health of application.
// If application unhealthy it must be restarted by external tools such as SIGINT
func (h *HealthChecksHandler) HandleHealthChecks(writer http.ResponseWriter, request *http.Request) {
	err := h.svc.CheckRepoHealth()
	if err != nil {
		log.Error(err)
		//writer.WriteHeader(http.StatusNotAcceptable)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
