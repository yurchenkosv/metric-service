package clients

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/yurchenkosv/metric-service/internal/model"
	"time"
)

type MetricServerClient struct {
	metricServerAddress string
	client              *resty.Client
}

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
