package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/types"
)

func TestCollectMemMetrics(t *testing.T) {
	pollCount := int64(1)
	type want struct {
		metrics       types.Metrics
		metricsLength int
	}
	tests := []struct {
		name      string
		pollCount int
		want      want
	}{
		{
			name:      "Should return map with metrics",
			pollCount: 1,
			want: want{
				metricsLength: 29,
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
			result := CollectMemMetrics(tt.pollCount)
			assert.IsType(t, types.Metrics{}, result)
			assert.NotEmpty(t, result.Metric)
			assert.Equal(t, len(result.Metric), tt.want.metricsLength)
			assert.Equal(t, int64(tt.pollCount), *result.Metric[28].Delta)
		})
	}
}
