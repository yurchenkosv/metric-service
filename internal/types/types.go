package types

import (
	"fmt"
	"time"
)

type Gauge float64
type Counter int64

type URLServer struct {
	host   string
	port   string
	schema string
}

func (s *URLServer) Build() string {
	if s.schema != "" {
		return fmt.Sprintf("%s://%s:%s", s.schema, s.host, s.port)
	}
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

func (s *URLServer) SetHost(host string) *URLServer {
	s.host = host
	return s
}

func (s *URLServer) SetPort(port string) *URLServer {
	s.port = port
	return s
}

func (s *URLServer) SetSchema(schema string) *URLServer {
	s.schema = schema
	return s
}

type Metrics struct {
	Metric []Metric
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}
