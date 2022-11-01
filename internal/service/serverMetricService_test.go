package service

import (
	"context"
	errors2 "errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	mock_repository "github.com/yurchenkosv/metric-service/internal/mockRepository"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/repository"
)

func TestServerMetricService_AddMetric(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, metric model.Metric, ctx context.Context)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		metric model.Metric
		ctx    context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		behavior mockBehavior
	}{
		{
			name:   "should successfuly add counter",
			fields: fields{config: &config.ServerConfig{}},
			args: args{
				metric: model.Metric{
					ID:    "testCounter",
					MType: "counter",
					Delta: model.NewCounter(15),
				},
				ctx: context.Background(),
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, metric model.Metric, ctx context.Context) {
				s.EXPECT().SaveCounter(ctx, metric.ID, *metric.Delta).Return(nil)
			},
		},
		{
			name:   "should successfuly add gauge",
			fields: fields{config: &config.ServerConfig{}},
			args: args{
				metric: model.Metric{
					ID:    "testGauge",
					MType: "gauge",
					Value: model.NewGauge(12.5),
				},
				ctx: context.Background(),
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, metric model.Metric, ctx context.Context) {
				s.EXPECT().SaveGauge(ctx, metric.ID, *metric.Value).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.metric, tt.args.ctx)
			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.AddMetric(tt.args.ctx, tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("AddMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerMetricService_AddMetricBatch(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, metrics model.Metrics, ctx context.Context)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		metrics model.Metrics
		ctx     context.Context
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		behavior mockBehavior
	}{
		{
			name:   "should successfully save metric batch",
			fields: fields{config: &config.ServerConfig{}},
			args: args{
				ctx: context.Background(),
				metrics: model.Metrics{Metric: []model.Metric{
					{
						ID:    "testGauge",
						MType: "gauge",
						Value: model.NewGauge(12.51233),
					},
					{
						ID:    "testCounter",
						MType: "counter",
						Delta: model.NewCounter(15),
					},
				}},
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, metrics model.Metrics, ctx context.Context) {
				s.EXPECT().SaveMetricsBatch(ctx, metrics.Metric).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.metrics, tt.args.ctx)
			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.AddMetricBatch(tt.args.ctx, tt.args.metrics); (err != nil) != tt.wantErr {
				t.Errorf("AddMetricBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerMetricService_CreateSignedHash(t *testing.T) {
	type fields struct {
		config *config.ServerConfig
		repo   repository.Repository
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		want          string
		wantErr       bool
		wantErrorType error
	}{
		{
			name: "shoud create correct hash with gauge",
			fields: fields{
				config: &config.ServerConfig{HashKey: "test"},
				repo:   repository.NewMapRepo(),
			},
			args: args{
				msg: "testGauge:gauge:12.5",
			},
			want:          "6d0b338da23630f6d0f3cd53d6f60e5140e91c39f346475f24b44544c79abafd",
			wantErr:       false,
			wantErrorType: nil,
		},
		{
			name: "shoud create correct hash with counter",
			fields: fields{
				config: &config.ServerConfig{HashKey: "test"},
				repo:   repository.NewMapRepo(),
			},
			args: args{
				msg: "testCounter:counter:7",
			},
			want:          "4338e1ebc35867905d5294b9e3c8b9196b3c54db1917e5db630e37c496956fc3",
			wantErr:       false,
			wantErrorType: nil,
		},
		{
			name: "shoud return NoEncryptionKeyFound error",
			fields: fields{
				config: &config.ServerConfig{},
				repo:   repository.NewMapRepo(),
			},
			args: args{
				msg: "testCounter:counter:7",
			},
			want:          "",
			wantErr:       true,
			wantErrorType: &errors.NoEncryptionKeyFoundError{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   tt.fields.repo,
			}
			got, err := s.CreateSignedHash(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSignedHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				assert.IsType(t, tt.wantErrorType, err)
			}
			if got != tt.want {
				t.Errorf("CreateSignedHash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerMetricService_GetAllMetrics(t *testing.T) {
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
		want     *model.Metrics
		wantErr  bool
		args     args
		behavior mockBehavior
	}{
		{
			name: "should return all metrics",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				ctx: context.Background(),
			},
			want: &model.Metrics{Metric: []model.Metric{
				{
					ID:    "testGauge",
					MType: "gauge",
					Value: model.NewGauge(12.51233),
				},
				{
					ID:    "testCounter",
					MType: "counter",
					Delta: model.NewCounter(15),
				},
			}},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, ctx context.Context) {
				metrics := model.Metrics{Metric: []model.Metric{
					{
						ID:    "testGauge",
						MType: "gauge",
						Value: model.NewGauge(12.51233),
					},
					{
						ID:    "testCounter",
						MType: "counter",
						Delta: model.NewCounter(15),
					},
				}}
				s.EXPECT().GetAllMetrics(ctx).Return(&metrics, nil)
			},
		},
		{
			name: "should return error",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args:    args{ctx: context.Background()},
			want:    nil,
			wantErr: true,
			behavior: func(s *mock_repository.MockRepository, ctx context.Context) {
				s.EXPECT().GetAllMetrics(ctx).Return(nil, errors2.New("testError"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.ctx)

			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			got, err := s.GetAllMetrics(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServerMetricService_GetMetricByKey(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, name string, metric *model.Metric, ctx context.Context)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		name string
		ctx  context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *model.Metric
		wantErr     bool
		wantErrType error
		behavior    mockBehavior
	}{
		{
			name: "should return counter",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				name: "testCounter",
				ctx:  context.Background(),
			},
			want: &model.Metric{
				ID:    "testCounter",
				MType: "counter",
				Delta: model.NewCounter(7),
				Value: nil,
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric, ctx context.Context) {
				s.EXPECT().GetMetricByKey(ctx, name).Return(metric, nil)
			},
		},
		{
			name: "should return gauge",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				name: "testGauge",
				ctx:  context.Background(),
			},
			want: &model.Metric{
				ID:    "testGauge",
				MType: "gauge",
				Delta: nil,
				Value: model.NewGauge(12.5),
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric, ctx context.Context) {
				s.EXPECT().GetMetricByKey(ctx, name).Return(metric, nil)
			},
		},
		{
			name: "should return error",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				name: "testGauge",
				ctx:  context.Background(),
			},
			want:        nil,
			wantErr:     true,
			wantErrType: errors.MetricNotFoundError{MetricName: "testGauge"},
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric, ctx context.Context) {
				s.EXPECT().GetMetricByKey(ctx, name).Return(nil, errors2.New("testError"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.name, tt.want, tt.args.ctx)

			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			got, err := s.GetMetricByKey(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetricByKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.wantErr {
				assert.IsType(t, &errors.MetricNotFoundError{MetricName: tt.args.name}, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMetricByKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestServerMetricService_LoadMetricsFromDisk(t *testing.T) {
//	type fields struct {
//		config *config.ServerConfig
//		repo   repository.Repository
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &ServerMetricService{
//				config: tt.fields.config,
//				repo:   tt.fields.repo,
//			}
//			if err := s.LoadMetricsFromDisk(); (err != nil) != tt.wantErr {
//				t.Errorf("LoadMetricsFromDisk() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}
//
//func TestServerMetricService_SaveMetricsToDisk(t *testing.T) {
//	type fields struct {
//		config *config.ServerConfig
//		repo   repository.Repository
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &ServerMetricService{
//				config: tt.fields.config,
//				repo:   tt.fields.repo,
//			}
//			if err := s.SaveMetricsToDisk(); (err != nil) != tt.wantErr {
//				t.Errorf("SaveMetricsToDisk() error = %v, wantErr %v", err, tt.wantErr)
//			}
//		})
//	}
//}

func TestNewServerMetricService(t *testing.T) {
	type args struct {
		cnf  *config.ServerConfig
		repo repository.Repository
	}
	tests := []struct {
		name string
		args args
		want *ServerMetricService
	}{
		{
			name: "should return ServerMetricService",
			args: args{
				cnf:  &config.ServerConfig{},
				repo: &repository.PostgresRepo{},
			},
			want: &ServerMetricService{
				config: &config.ServerConfig{},
				repo:   &repository.PostgresRepo{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewServerMetricService(tt.args.cnf, tt.args.repo), "NewServerMetricService(%v, %v)", tt.args.cnf, tt.args.repo)
		})
	}
}
