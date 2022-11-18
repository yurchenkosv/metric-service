package clients

import (
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// MetricServerClient struct with address and resty client to perform http requests.
type MetricServerClient struct {
	client              *resty.Client
	metricServerAddress string
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
func (c *MetricServerClient) PushMetrics(metrics string) {
	go c.pushToServer(metrics)
}

func (c *MetricServerClient) pushToServer(msg string) {
	_, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(msg).
		Post("/updates")
	if err != nil {
		log.Error(err)
	}
}
