package service

import (
	errors2 "errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/errors"
	mock_repository "github.com/yurchenkosv/metric-service/internal/mocks"
	"github.com/yurchenkosv/metric-service/internal/model"
	"github.com/yurchenkosv/metric-service/internal/repository"
	"reflect"
	"testing"
)

func TestServerMetricService_AddMetric(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, metric model.Metric)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		metric model.Metric
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
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, metric model.Metric) {
				s.EXPECT().SaveCounter(metric.ID, *metric.Delta).Return(nil)
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
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, metric model.Metric) {
				s.EXPECT().SaveGauge(metric.ID, *metric.Value).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.metric)
			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.AddMetric(tt.args.metric); (err != nil) != tt.wantErr {
				t.Errorf("AddMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServerMetricService_AddMetricBatch(t *testing.T) {
	type mockBehavior func(s *mock_repository.MockRepository, metrics model.Metrics)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		metrics model.Metrics
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
			behavior: func(s *mock_repository.MockRepository, metrics model.Metrics) {
				s.EXPECT().SaveMetricsBatch(metrics.Metric).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.metrics)
			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			if err := s.AddMetricBatch(tt.args.metrics); (err != nil) != tt.wantErr {
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
	type mockBehavior func(s *mock_repository.MockRepository)
	type fields struct {
		config *config.ServerConfig
	}
	tests := []struct {
		name     string
		fields   fields
		want     *model.Metrics
		wantErr  bool
		behavior mockBehavior
	}{
		{
			name: "should return all metrics",
			fields: fields{
				config: &config.ServerConfig{},
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
			behavior: func(s *mock_repository.MockRepository) {
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
				s.EXPECT().GetAllMetrics().Return(&metrics, nil)
			},
		},
		{
			name: "should return error",
			fields: fields{
				config: &config.ServerConfig{},
			},
			want:    nil,
			wantErr: true,
			behavior: func(s *mock_repository.MockRepository) {
				s.EXPECT().GetAllMetrics().Return(nil, errors2.New("testError"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo)

			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			got, err := s.GetAllMetrics()
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
	type mockBehavior func(s *mock_repository.MockRepository, name string, metric *model.Metric)
	type fields struct {
		config *config.ServerConfig
	}
	type args struct {
		name string
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
			},
			want: &model.Metric{
				ID:    "testCounter",
				MType: "counter",
				Delta: model.NewCounter(7),
				Value: nil,
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric) {
				s.EXPECT().GetMetricByKey(name).Return(metric, nil)
			},
		},
		{
			name: "should return gauge",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				name: "testGauge",
			},
			want: &model.Metric{
				ID:    "testGauge",
				MType: "gauge",
				Delta: nil,
				Value: model.NewGauge(12.5),
			},
			wantErr: false,
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric) {
				s.EXPECT().GetMetricByKey(name).Return(metric, nil)
			},
		},
		{
			name: "should return error",
			fields: fields{
				config: &config.ServerConfig{},
			},
			args: args{
				name: "testGauge",
			},
			want:        nil,
			wantErr:     true,
			wantErrType: errors.MetricNotFoundError{MetricName: "testGauge"},
			behavior: func(s *mock_repository.MockRepository, name string, metric *model.Metric) {
				s.EXPECT().GetMetricByKey(name).Return(nil, errors2.New("testError"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mock_repository.NewMockRepository(ctrl)
			tt.behavior(repo, tt.args.name, tt.want)

			s := &ServerMetricService{
				config: tt.fields.config,
				repo:   repo,
			}
			got, err := s.GetMetricByKey(tt.args.name)
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
