package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/yurchenkosv/metric-service/internal/config"
	mock_repository "github.com/yurchenkosv/metric-service/internal/mockRepository"
)

func TestHealthCheckService_CheckRepoHealth(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository)
	type fields struct {
		config *config.ServerConfig
	}
	tests := []struct {
		name     string
		fields   fields
		wantErr  bool
		behavior mockBehavior
	}{
		{
			name: "should successfully check health",
			behavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().Ping().Return(nil)
			},
			wantErr: false,
			fields:  struct{ config *config.ServerConfig }{config: &config.ServerConfig{}},
		},
		{
			name: "should return error",
			behavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().Ping().Return(errors.New("no healthy storage"))
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
			tt.behavior(repo)
			s := HealthCheckService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.CheckRepoHealth(); (err != nil) != tt.wantErr {
				t.Errorf("CheckRepoHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
