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
	err := h.svc.CheckRepoHealth(ctx)
	if err != nil {
		return &api.Pong{
			Status: api.HealthStatus_error,
		}, err
	}
	return &api.Pong{
		Status: api.HealthStatus_healthy,
	}, nil
}
