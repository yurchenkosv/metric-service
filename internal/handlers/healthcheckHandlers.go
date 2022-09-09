package handlers

import (
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/service"
	"net/http"
)

type HealthChecksHandler struct {
	svc *service.HealthCheckService
}

func NewHealthCheckHandler(svc *service.HealthCheckService) *HealthChecksHandler {
	return &HealthChecksHandler{
		svc: svc,
	}
}

func (h *HealthChecksHandler) HandleHealthChecks(writer http.ResponseWriter, request *http.Request) {
	err := h.svc.CheckRepoHealth()
	if err != nil {
		log.Error(err)
		//writer.WriteHeader(http.StatusNotAcceptable)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
