package service

import (
	"context"
	"encoding/json"
	errors2 "errors"
	"fmt"
	"os"
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
	type mockBehavior func(ctx context.Context, s *mock_repository.MockRepository, metric model.Metric)
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
			behavior: func(ctx context.Context, s *mock_repository.MockRepository, metric model.Metric) {
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
			behavior: func(ctx context.Context, s *mock_repository.MockRepository, metric model.Metric) {
				s.EXPECT().SaveGauge(ctx, metric.ID, *metric.Value).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(tt.args.ctx, repo, tt.args.metric)
			s := &ServerMetricService{
				config:            tt.fields.config,
				saveMetricsToDisk: false,
				repo:              repo,
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

func TestServerMetricService_LoadMetricsFromDisk(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_repository.MockRepository, metrics model.Metrics)
	type fields struct {
		config *config.ServerConfig
		repo   repository.Repository
	}
	type args struct {
		ctx    context.Context
		metric model.Metric
	}

	tests := []struct {
		name     string
		fields   fields
		wantErr  assert.ErrorAssertionFunc
		behavior mockBehavior
		before   func(fileLocation string, metrics model.Metrics)
		args     args
	}{
		{
			name: "should successfully load metrics",
			fields: fields{
				config: &config.ServerConfig{
					StoreFile: "./metric_save",
				},
				repo: nil,
			},
			wantErr: assert.NoError,
			args: args{
				ctx: context.Background(),
				metric: model.Metric{
					ID:    "testGauge",
					MType: "gauge",
					Value: model.NewGauge(0.25),
				},
			},
			before: func(fileLocation string, metrics model.Metrics) {
				fileBits := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				file, _ := os.OpenFile(fileLocation, fileBits, 0600)
				data, _ := json.Marshal(metrics)
				file.Write(data)
				file.Close()
			},
			behavior: func(ctx context.Context, s *mock_repository.MockRepository, metrics model.Metrics) {
				s.EXPECT().SaveMetricsBatch(ctx, metrics.Metric).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			metrics := model.Metrics{[]model.Metric{tt.args.metric}}
			tt.behavior(tt.args.ctx, repo, metrics)
			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			tt.before(tt.fields.config.StoreFile, metrics)
			tt.wantErr(t, s.LoadMetricsFromDisk(tt.args.ctx), fmt.Sprintf("LoadMetricsFromDisk()"))
		})
	}
}

func TestServerMetricService_SaveMetricsToDisk1(t *testing.T) {
	type mockBehavior func(ctx context.Context, s *mock_repository.MockRepository, metrics model.Metrics)
	type fields struct {
		config            *config.ServerConfig
		repo              repository.Repository
		saveMetricsToDisk bool
	}
	type args struct {
		ctx     context.Context
		metrics model.Metrics
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		behavior mockBehavior
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "should successfuly save metrics to disk",
			fields: fields{
				config: &config.ServerConfig{
					StoreFile: "./metric_save",
				},
				saveMetricsToDisk: true,
			},
			args: args{
				ctx: context.Background(),
				metrics: struct{ Metric []model.Metric }{Metric: []model.Metric{
					{
						ID:    "TestGauge",
						MType: "gauge",
						Value: model.NewGauge(5.2),
					},
				}},
			},
			behavior: func(ctx context.Context, s *mock_repository.MockRepository, metrics model.Metrics) {
				s.EXPECT().GetAllMetrics(ctx).Return(&metrics, nil)
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(tt.args.ctx, repo, tt.args.metrics)

			s := &ServerMetricService{
				config:            tt.fields.config,
				repo:              repo,
				saveMetricsToDisk: tt.fields.saveMetricsToDisk,
			}

			tt.wantErr(t, s.SaveMetricsToDisk(tt.args.ctx), fmt.Sprintf("SaveMetricsToDisk(%v)", tt.args.ctx))
		})
	}
}

func TestServerMetricService_Shutdown(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository)
	type fields struct {
		config            *config.ServerConfig
		repo              repository.Repository
		saveMetricsToDisk bool
	}
	tests := []struct {
		name     string
		fields   fields
		behavior mockBehavior
	}{
		{
			name: "should successfully shutdown service",
			fields: fields{
				config:            &config.ServerConfig{},
				repo:              nil,
				saveMetricsToDisk: false,
			},
			behavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().Shutdown()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo)
			s := ServerMetricService{
				config:            tt.fields.config,
				repo:              repo,
				saveMetricsToDisk: tt.fields.saveMetricsToDisk,
			}
			s.Shutdown()
		})
	}
}
