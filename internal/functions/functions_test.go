package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/types"
)

func TestCollectMemMetrics(t *testing.T) {
	pollCount := int64(1)
	type want struct {
		metrics types.Metrics
	}
	tests := []struct {
		name      string
		pollCount int
		cfg       types.AgentConfig
		want      want
	}{
		{
			name:      "Should return map with metrics",
			pollCount: 1,
			cfg:       types.AgentConfig{},
			want: want{
				metrics: types.Metrics{Metric: []types.Metric{
					{
						ID:    "PollCount",
						MType: "counter",
						Delta: &pollCount,
						Value: nil,
					},
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CollectMetrics(tt.pollCount, &tt.cfg)
			assert.IsType(t, types.Metrics{}, result)
			assert.NotEmpty(t, result.Metric)
		})
	}
}
