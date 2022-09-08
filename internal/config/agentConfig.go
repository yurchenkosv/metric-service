package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"time"
)

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	HashKey        string        `env:"KEY"`
}

func (c *AgentConfig) Parse() error {
	flag.StringVar(&c.Address, "a", "localhost:8080", "http address to send metrics in format localhost:8080")
	flag.DurationVar(&c.ReportInterval, "r", 10*time.Second, "interval to send metrics to server. Inactive for server.")
	flag.DurationVar(&c.PollInterval, "p", 2*time.Second, "Interval to collect metrics. Inactive for server.")
	flag.StringVar(&c.HashKey, "k", "", "key to create hash")
	flag.Parse()

	err := env.Parse(c)
	return err
}
