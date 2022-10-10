//asdasdasd
package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

// ServerConfig struct with fields, useful for configuring metrics server.
type ServerConfig struct {
	StoreInterval time.Duration `env:"STORE_INTERVAL"` // when to flush metrics to disk.
	Address       string        `env:"ADDRESS"`        // server address to bind to.
	StoreFile     string        `env:"STORE_FILE"`     // path to file where metrics are stored.
	Restore       bool          `env:"RESTORE"`        // If set to true, read StoreFile to restore metrics state
	HashKey       string        `env:"KEY"`            // key to create/validate hash
	DBDsn         string        `env:"DATABASE_DSN"`   // Postgres connection string in DSN format
}

// NewServerConfig constructor returns pointer to new ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

// Parse method to fulfill ServerConfig fields.
// It reads flags and env variables.
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
