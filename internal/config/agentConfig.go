package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// AgentConfig struct with fields, useful for configuring metrics agent.
type AgentConfig struct {
	Address        string        `env:"ADDRESS"`         // http address of metric server
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // interval to send metrics to server
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // interval to collect metrics
	HashKey        string        `env:"KEY"`             // key to create hash
}

// Parse method to fulfill AgentConfig fields.
// It reads flags and env variables.
func (c *AgentConfig) Parse() error {
	flag.StringVar(&c.Address, "a", "localhost:8080", "http address to send metrics in format localhost:8080")
	flag.DurationVar(&c.ReportInterval, "r", 10*time.Second, "interval to send metrics to server. Inactive for server.")
	flag.DurationVar(&c.PollInterval, "p", 2*time.Second, "Interval to collect metrics. Inactive for server.")
	flag.StringVar(&c.HashKey, "k", "", "key to create hash")
	flag.Parse()

	err := env.Parse(c)
	return err
}
