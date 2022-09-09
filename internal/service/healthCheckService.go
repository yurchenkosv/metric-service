package service

import (
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	"github.com/yurchenkosv/metric-service/internal/repository"
)

type HealthCheckService struct {
	config *config.ServerConfig
	repo   repository.Repository
}

func NewHealthCheckService(cnf *config.ServerConfig, repo repository.Repository) *HealthCheckService {
	return &HealthCheckService{
		config: cnf,
		repo:   repo,
	}
}

func (s HealthCheckService) CheckRepoHealth() error {
	if s.repo.Ping() != nil {
		return &errors.HealthCheckError{HealthcheckType: "Repository"}
	}
	return nil
}
