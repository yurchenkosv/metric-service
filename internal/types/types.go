package types

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"io"
	"net/http"
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
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	Key            string        `env:"KEY"`
}

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	Key           string        `env:"KEY"`
	DBDsn         string        `env:"DATABASE_DSN"`
}

func (c *AgentConfig) Parse() error {
	flag.StringVar(&c.Address, "a", "localhost:8080", "http address to send metrics in format localhost:8080")
	flag.DurationVar(&c.ReportInterval, "r", 10*time.Second, "interval to send metrics to server. Inactive for server.")
	flag.DurationVar(&c.PollInterval, "p", 2*time.Second, "Interval to collect metrics. Inactive for server.")
	flag.StringVar(&c.Key, "k", "", "key to create hash")
	flag.Parse()

	err := env.Parse(c)
	return err
}

func (c *ServerConfig) Parse() error {
	flag.StringVar(&c.Address, "a", "localhost:8080", "http address in format localhost:8080")
	flag.DurationVar(&c.StoreInterval, "i", 300*time.Second, "when to flush metrics to disk. Inactive for agent.")
	flag.StringVar(&c.StoreFile, "f", "/tmp/devops-metrics-db.json", "path to file where metrics are stored. Inactive for agent.")
	flag.BoolVar(&c.Restore, "r", true, "If set to true, read file in -f flag to restore metrics state")
	flag.StringVar(&c.Key, "k", "", "key to create/validate hash")
	flag.StringVar(&c.DBDsn, "d", "", "")
	flag.Parse()

	err := env.Parse(c)
	if c.DBDsn != "" {
		c.StoreFile = ""
		c.Restore = false
	}
	return err
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}
