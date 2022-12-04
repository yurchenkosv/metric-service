package handlers

import (
	"context"
	"github.com/yurchenkosv/metric-service/internal/api"
	"github.com/yurchenkosv/metric-service/internal/service"
)

type GRPCHealthChecksHandler struct {
	svc *service.HealthCheckService
}

// NewGRPCHealthCheckHandler sets service.HealthCheckService and returns pointer to HealthChecksHandler.
func NewGRPCHealthCheckHandler(svc *service.HealthCheckService) *GRPCHealthChecksHandler {
	return &GRPCHealthChecksHandler{
		svc: svc,
	}
}

func (h *GRPCHealthChecksHandler) GetHealthStatus(ctx context.Context, ping *api.Ping) (*api.Pong, error) {
	var pong *api.Pong
	err := h.svc.CheckRepoHealth(ctx)
	if err != nil {
		pong.Status = api.HealthStatus_error
		return pong, err
	}
	pong.Status = api.HealthStatus_healthy
	return pong, nil
}
