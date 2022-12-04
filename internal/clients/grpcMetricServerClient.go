package clients

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/api"
	"github.com/yurchenkosv/metric-service/internal/model"
	"google.golang.org/grpc"
)

// GRPCMetricServerClient struct with address and resty client to perform http requests.
type GRPCMetricServerClient struct {
	connect *grpc.ClientConn
	opts    []grpc.CallOption
}

// NewGRPCMetricServerClient constructor creates resty client, configures it.
// Returns pointer to MetricServerClient
func NewGRPCMetricServerClient(connect *grpc.ClientConn, opts ...grpc.CallOption) *GRPCMetricServerClient {
	return &GRPCMetricServerClient{
		connect: connect,
		opts:    opts,
	}
}

func (c *GRPCMetricServerClient) PushMetrics(metrics []model.Metric) {
	apiMetrics := &api.Metrics{}

	ctx := context.Background()
	for _, metric := range metrics {
		apiMetric, err := api.MetricToApiMetric(metric)
		if err != nil {
			log.Errorf("cannot transform metric %v to API metric. Error: %s", metric, err)
		}
		apiMetrics.Metrics = append(apiMetrics.Metrics, apiMetric)
	}

	go func(metrics *api.Metrics) {
		client := api.NewMetricServiceClient(c.connect)
		response, err := client.SaveMetrics(ctx, apiMetrics, c.opts...)
		if err != nil {
			log.Error(err)
		}
		log.Debug("send request, got response with grpc: ", response)
	}(apiMetrics)
}
