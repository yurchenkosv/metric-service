package clients

import (
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/model"
)

// MetricServerClient struct with address and resty client to perform http requests.
type MetricServerClient struct {
	metricServerAddress string
	client              *resty.Client
}

// NewMetricServerClient constructor creates resty client, configures it.
// Returns pointer to MetricServerClient
func NewMetricServerClient(address string) *MetricServerClient {
	client := &MetricServerClient{
		metricServerAddress: address,
		client:              resty.New(),
	}
	client.client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		SetBaseURL("http://" + client.metricServerAddress)
	return client
}

// PushMetrics method sends metrics to metric server in multiple threads via http.
func (c MetricServerClient) PushMetrics(metrics model.Metrics) {
	go func() {
		if len(metrics.Metric) > 0 {
			_, err := c.client.R().
				SetHeader("Content-Type", "application/json").
				SetBody(metrics.Metric).
				Post("/updates")
			if err != nil {
				log.Error(err)
			}
		}
	}()
}
