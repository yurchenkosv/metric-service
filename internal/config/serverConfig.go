package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	HashKey       string        `env:"KEY"`
	DBDsn         string        `env:"DATABASE_DSN"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) Parse() error {
	flag.StringVar(&c.Address, "a", "localhost:8080", "http address in format localhost:8080")
	flag.DurationVar(&c.StoreInterval, "i", 300*time.Second, "when to flush metrics to disk. Inactive for agent.")
	flag.StringVar(&c.StoreFile, "f", "/tmp/devops-metrics-db.json", "path to file where metrics are stored. Inactive for agent.")
	flag.BoolVar(&c.Restore, "r", true, "If set to true, read file in -f flag to restore metrics state")
	flag.StringVar(&c.HashKey, "k", "", "key to create/validate hash")
	flag.StringVar(&c.DBDsn, "d", "", "Postgres connection string")
	flag.Parse()

	err := env.Parse(c)
	if c.DBDsn != "" {
		c.Restore = false
		c.StoreFile = ""
	}
	return err
}
