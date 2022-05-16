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
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
}

func (c *AgentConfig) Parse() error {
	flag.StringVarP(&c.Address, "address", "a", c.Address, "http address to send metrics in format localhost:8080")
	flag.DurationVarP(&c.ReportInterval, "report", "r", c.ReportInterval, "interval to send metrics to server. Inactive for server.")
	flag.DurationVarP(&c.PollInterval, "poll", "p", c.PollInterval, "Interval to collect metrics. Inactive for server.")
	flag.Parse()

	err := env.Parse(c)
	return err
}

func (c *ServerConfig) Parse() error {
	flag.StringVarP(&c.Address, "address", "a", "localhost:8080", "http address in format localhost:8080")
	flag.DurationVarP(&c.StoreInterval, "interval", "i", 300*time.Second, "when to flush metrics to disk. Inactive for agent.")
	flag.StringVarP(&c.StoreFile, "filepath", "f", "/tmp/devops-metrics-db.json", "path to file where metrics are stored. Inactive for agent.")
	flag.BoolVarP(&c.Restore, "restore", "r", true, "If set to true, read file in -f flag to restore metrics state")
	flag.Parse()

	err := env.Parse(c)
	return err
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}
