package types

import (
	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
	"time"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Metric []Metric
}

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
}

type ServerConfig struct {
	Address       string        `env:"ADDRESS" envDefault:"localhost:8080"`
	StoreInterval time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile     string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore       bool          `env:"RESTORE" envDefault:"true"`
}

func (c *AgentConfig) Parse() error {
	err := env.Parse(c)
	flag.StringVar(&c.Address, "a", c.Address, "http address to send metrics in format localhost:8080")
	flag.DurationVar(&c.ReportInterval, "r", c.ReportInterval, "interval to send metrics to server. Inactive for server.")
	flag.DurationVar(&c.PollInterval, "p", c.PollInterval, "Interval to collect metrics. Inactive for server.")
	flag.Parse()

	return err
}

func (c *ServerConfig) Parse() error {
	err := env.Parse(c)
	flag.StringVar(&c.Address, "a", c.Address, "http address in format localhost:8080")
	flag.DurationVar(&c.StoreInterval, "i", c.StoreInterval, "when to flush metrics to disk. Inactive for agent.")
	flag.StringVar(&c.StoreFile, "f", c.StoreFile, "path to file where metrics are stored. Inactive for agent.")
	flag.BoolVar(&c.Restore, "r", c.Restore, "If set to true, read file in -f flag to restore metrics state")
	flag.Parse()

	return err
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}
