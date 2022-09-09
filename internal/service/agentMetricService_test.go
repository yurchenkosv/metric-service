package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/config"
	"github.com/yurchenkosv/metric-service/internal/model"
	"testing"
)

func TestAgentMetricService_CollectMetrics(t *testing.T) {
	type fields struct {
		config *config.AgentConfig
	}
	type args struct {
		poolCount int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   model.Metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAgentMetricService(tt.fields.config)
			assert.Equalf(t, tt.want, s.CollectMetrics(tt.args.poolCount), "CollectMetrics(%v)", tt.args.poolCount)
		})
	}
}

func TestAgentMetricService_CreateSignedHash(t *testing.T) {
	type fields struct {
		config *config.AgentConfig
	}
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "shoud create correct hash with gauge",
			fields: fields{
				config: &config.AgentConfig{HashKey: "test"},
			},
			args: args{
				msg: "testGauge:gauge:12.5",
			},
			want:    "6d0b338da23630f6d0f3cd53d6f60e5140e91c39f346475f24b44544c79abafd",
			wantErr: assert.NoError,
		},
		{
			name: "shoud create correct hash with counter",
			fields: fields{
				config: &config.AgentConfig{HashKey: "test"},
			},
			args: args{
				msg: "testCounter:counter:7",
			},
			want:    "4338e1ebc35867905d5294b9e3c8b9196b3c54db1917e5db630e37c496956fc3",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAgentMetricService(tt.fields.config)
			got, err := s.CreateSignedHash(tt.args.msg)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateSignedHash(%v)", tt.args.msg)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreateSignedHash(%v)", tt.args.msg)
		})
	}
}

func TestAgentMetricService_appendCounterMetric(t *testing.T) {
	type fields struct {
		config *config.AgentConfig
	}
	type args struct {
		name    string
		value   int64
		metrics *model.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should sucess add counter to metrics",
			fields: fields{
				config: &config.AgentConfig{},
			},
			args: args{
				name:    "",
				value:   7,
				metrics: &model.Metrics{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAgentMetricService(tt.fields.config)
			s.appendCounterMetric(tt.args.name, tt.args.value, tt.args.metrics)
		})
	}
}

func TestAgentMetricService_appendGaugeMetric(t *testing.T) {
	type fields struct {
		config *config.AgentConfig
	}
	type args struct {
		name    string
		value   float64
		metrics *model.Metrics
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should success add gauge to metrics",
			fields: fields{
				config: &config.AgentConfig{},
			},
			args: args{
				name:    "testGauge",
				value:   12.5,
				metrics: &model.Metrics{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewAgentMetricService(tt.fields.config)
			s.appendGaugeMetric(tt.args.name, tt.args.value, tt.args.metrics)
		})
	}
}

func TestNewAgentMetricService(t *testing.T) {
	type args struct {
		cfg *config.AgentConfig
	}
	tests := []struct {
		name string
		args args
		want *AgentMetricService
	}{
		{
			args: args{cfg: &config.AgentConfig{}},
			name: "should create AgentMetricService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want = &AgentMetricService{config: tt.args.cfg}
			assert.Equalf(t, tt.want, NewAgentMetricService(tt.args.cfg), "NewAgentMetricService(%v)", tt.args.cfg)
		})
	}
}