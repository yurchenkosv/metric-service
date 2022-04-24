package functions

import (
	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/types"
	"testing"
)

func TestCollectMemMetrics(t *testing.T) {
	type want struct {
		metrics       types.MemMetrics
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
				metrics:       types.MemMetrics{PollCount: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CollectMemMetrics(tt.poolCount)
			assert.IsType(t, types.MemMetrics{}, result)
			assert.NotEmpty(t, result.GaugeMetrics)
			assert.Equal(t, len(result.GaugeMetrics), tt.want.metricsLength)
			assert.Equal(t, tt.want.metrics.PollCount, result.PollCount)
		})
	}
}
