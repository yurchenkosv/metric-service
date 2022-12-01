package clients

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/model"
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

func (c *MetricServerClient) SetScheme(scheme string) *MetricServerClient {
	c.client.
		SetBaseURL(scheme + "://" + c.metricServerAddress)
	return c
}

func (c *MetricServerClient) SetHeader(key, val string) *MetricServerClient {
	c.client.SetHeader(key, val)
	return c
}

// PushMetrics method sends metrics to metric server in multiple threads via http.
func (c *MetricServerClient) PushMetrics(metrics []model.Metric) {
	go c.pushToServer(metrics)
}

func (c *MetricServerClient) WithTLS(tlsConfig *tls.Config) *MetricServerClient {
	c.client.SetTLSClientConfig(tlsConfig)
	return c
}

func (c *MetricServerClient) pushToServer(metrics []model.Metric) {
	_, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post("/updates")
	if err != nil {
		log.Error(err)
	}
}
