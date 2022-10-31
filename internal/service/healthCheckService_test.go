package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yurchenkosv/metric-service/internal/config"
	mock_repository "github.com/yurchenkosv/metric-service/internal/mockRepository"
	"github.com/yurchenkosv/metric-service/internal/repository"
)

func TestHealthCheckService_CheckRepoHealth(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, ctx context.Context)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		args     args
		behavior mockBehavior
	}{
		{
			name: "should successfully check health",
			behavior: func(s *mock_repository.MockRepository, ctx context.Context) {
				s.EXPECT().Ping(ctx).Return(nil)
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			fields:  struct{ config *config.ServerConfig }{config: &config.ServerConfig{}},
		},
		{
			name: "should return error",
			behavior: func(s *mock_repository.MockRepository, ctx context.Context) {
				s.EXPECT().Ping(ctx).Return(errors.New("no healthy storage"))
			},
			wantErr: true,
			fields:  struct{ config *config.ServerConfig }{config: &config.ServerConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.ctx)
			s := HealthCheckService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.CheckRepoHealth(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("CheckRepoHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewHealthCheckService(t *testing.T) {
	type args struct {
		cnf  *config.ServerConfig
		repo *repository.PostgresRepo
	}
	tests := []struct {
		name string
		args args
		want *HealthCheckService
	}{
		{
			name: "chould successfully create healthcheck service",
			args: args{
				cnf:  &config.ServerConfig{},
				repo: &repository.PostgresRepo{},
			},
			want: &HealthCheckService{
				config: &config.ServerConfig{},
				repo:   &repository.PostgresRepo{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewHealthCheckService(tt.args.cnf, tt.args.repo), "NewHealthCheckService(%v, %v)", tt.args.cnf, tt.args.repo)
		})
	}
}
