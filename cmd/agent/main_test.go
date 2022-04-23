package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_collectMemMetrics(t *testing.T) {
	type want struct {
		metrics       MemMetrics
		metricsLength int
	}
	tests := []struct {
		name      string
		poolCount int
		want      want
	}{
		{
			name:      "Should return map with metrics",
			poolCount: 1,
			want: want{
				metricsLength: 28,
				metrics:       MemMetrics{PollCount: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := collectMemMetrics(tt.poolCount)
			assert.IsType(t, MemMetrics{}, result)
			assert.NotEmpty(t, result.gaugeMetrics)
			assert.Equal(t, len(result.gaugeMetrics), tt.want.metricsLength)
			assert.Equal(t, tt.want.metrics.PollCount, result.PollCount)
		})
	}
}
