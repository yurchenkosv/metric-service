package service

import (
	"context"

	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/repository"
)

// HealthCheckService is for checking health of metrics server.
type HealthCheckService struct {
	config *config.ServerConfig
	repo   repository.Repository
}

// NewHealthCheckService create new NewHealthCheckService with filled fields and returns pointer on this object.
func NewHealthCheckService(cnf *config.ServerConfig, repo repository.Repository) *HealthCheckService {
	return &HealthCheckService{
		config: cnf,
		repo:   repo,
	}
}

// CheckRepoHealth is main method of NewHealthCheckService. It checks that connection to DB could be made.
// Returns error if not.
func (s HealthCheckService) CheckRepoHealth(ctx context.Context) error {
	if s.repo.Ping(ctx) != nil {
		return &errors.HealthCheckError{HealthcheckType: "Repository"}
	}
	return nil
}
