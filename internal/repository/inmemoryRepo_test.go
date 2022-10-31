package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yurchenkosv/metric-service/internal/model"
)

func TestNewMapRepo(t *testing.T) {
	tests := []struct {
		name string
		want *mapStorage
	}{
		{
			name: "should create map repo",
			want: NewMapRepo(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMapRepo(), "NewMapRepo()")
		})
	}
}

func Test_mapStorage_GetAllMetrics(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    *model.Metrics
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully get metrics",
			fields: fields{
				GaugeMetric: map[string]model.Gauge{
					"RandomGauge": model.Gauge(500.20),
				},
				CounterMetric: map[string]model.Counter{
					"RandomCounter": model.Counter(100),
				},
			},
			args: args{ctx: context.Background()},
			want: &model.Metrics{Metric: []model.Metric{
				{
					ID:    "RandomCounter",
					MType: "counter",
					Delta: model.NewCounter(100),
				},
				{
					ID:    "RandomGauge",
					MType: "gauge",
					Value: model.NewGauge(500.20),
				},
			}},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			got, err := m.GetAllMetrics(tt.args.ctx)
			if !tt.wantErr(t, err, "GetAllMetrics()") {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllMetrics()")
		})
	}
}

func Test_mapStorage_GetMetricByKey(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		key string
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Metric
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should successfully return counter by key",
			fields: fields{
				GaugeMetric: map[string]model.Gauge{},
				CounterMetric: map[string]model.Counter{
					"RandomCounter": 100,
				},
			},
			args: args{
				key: "RandomCounter",
				ctx: context.Background(),
			},
			want: &model.Metric{
				ID:    "RandomCounter",
				MType: "counter",
				Delta: model.NewCounter(100),
			},
			wantErr: assert.NoError,
		},
		{
			name: "should successfully return gauge by key",
			fields: fields{
				GaugeMetric: map[string]model.Gauge{
					"RandomGauge": model.Gauge(500.24),
				},
				CounterMetric: map[string]model.Counter{},
			},
			args: args{
				key: "RandomGauge",
				ctx: context.Background(),
			},
			want: &model.Metric{
				ID:    "RandomGauge",
				MType: "gauge",
				Value: model.NewGauge(500.24),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			got, err := m.GetMetricByKey(tt.args.key, tt.args.ctx)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMetricByKey(%v)", tt.args.key)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMetricByKey(%v)", tt.args.key)
		})
	}
}

func Test_mapStorage_Migrate(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should execute migrate method",
			fields: fields{
				GaugeMetric:   nil,
				CounterMetric: nil,
			},
			args: args{
				path: "/fakepath",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			m.Migrate(tt.args.path)
		})
	}
}

func Test_mapStorage_Ping(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should execute ping method",
			fields: fields{
				GaugeMetric:   nil,
				CounterMetric: nil,
			},
			args:    args{ctx: context.Background()},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			tt.wantErr(t, m.Ping(tt.args.ctx), "Ping()")
		})
	}
}

func Test_mapStorage_SaveCounter(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		name string
		val  model.Counter
		ctx  context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Counter
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should save counter to map",
			fields: fields{
				GaugeMetric:   map[string]model.Gauge{},
				CounterMetric: map[string]model.Counter{},
			},
			args: args{
				name: "RandomCounter",
				val:  212,
				ctx:  context.Background(),
			},
			want:    model.Counter(212),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}

			tt.wantErr(t, m.SaveCounter(tt.args.name, tt.args.val, tt.args.ctx), fmt.Sprintf("SaveCounter(%v, %v)", tt.args.name, tt.args.val))
			if val, ok := m.CounterMetric[tt.name]; ok {
				got := val
				assert.Equalf(t, tt.want, got, "GetMetricByKey(%v)", tt.args.name)
			}
		})
	}
}

func Test_mapStorage_SaveGauge(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		name string
		val  model.Gauge
		ctx  context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		want    model.Gauge
	}{
		{
			name: "should succesfully save gauge metric",
			fields: fields{
				GaugeMetric:   map[string]model.Gauge{},
				CounterMetric: map[string]model.Counter{},
			},
			args: args{
				name: "RandomGauge",
				val:  500.24,
				ctx:  context.Background(),
			},
			wantErr: assert.NoError,
			want:    model.Gauge(500.24),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			tt.wantErr(t, m.SaveGauge(tt.args.name, tt.args.val, tt.args.ctx), fmt.Sprintf("SaveGauge(%v, %v)", tt.args.name, tt.args.val))
			if val, ok := m.GaugeMetric[tt.name]; ok {
				got := val
				assert.Equalf(t, tt.want, got, "GetMetricByKey(%v)", tt.args.name)
			}
		})
	}
}

func Test_mapStorage_SaveMetricsBatch(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	type args struct {
		metrics []model.Metric
		ctx     context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
		want    *model.Metrics
	}{
		{
			name: "should successfully save metrics batch",
			fields: fields{
				GaugeMetric:   map[string]model.Gauge{},
				CounterMetric: map[string]model.Counter{},
			},
			args: args{
				ctx: context.Background(),
				metrics: []model.Metric{
					{
						ID:    "RandomCounter",
						MType: "counter",
						Delta: model.NewCounter(123),
					},
					{
						ID:    "RandomGauge",
						MType: "gauge",
						Value: model.NewGauge(500.24),
					},
				},
			},
			wantErr: assert.NoError,
			want: &model.Metrics{Metric: []model.Metric{
				{
					ID:    "RandomCounter",
					MType: "counter",
					Delta: model.NewCounter(123),
				},
				{
					ID:    "RandomGauge",
					MType: "gauge",
					Value: model.NewGauge(500.24),
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			tt.wantErr(t, m.SaveMetricsBatch(tt.args.metrics, tt.args.ctx), fmt.Sprintf("SaveMetricsBatch(%v)", tt.args.metrics))
			got, _ := m.GetAllMetrics(tt.args.ctx)
			assert.Equalf(t, tt.want, got, "GetAllMetrics()")
		})
	}
}

func Test_mapStorage_Shutdown(t *testing.T) {
	type fields struct {
		GaugeMetric   map[string]model.Gauge
		CounterMetric map[string]model.Counter
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "should successfully call shutdown method",
			fields: fields{
				GaugeMetric:   nil,
				CounterMetric: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mapStorage{
				GaugeMetric:   tt.fields.GaugeMetric,
				CounterMetric: tt.fields.CounterMetric,
			}
			m.Shutdown()
		})
	}
}
