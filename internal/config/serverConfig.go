package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

const serverDefaultStoreInterval time.Duration = 300 * time.Second
const serverDefaultAddress string = "localhost:8080"
const serverDefaultStoreFile string = "/tmp/devops-metrics-db.json"
const serverDefaultRestoreCondition bool = true
const serverDefaultHashKey string = ""
const serverDefaultDBDsn string = ""
const serverDefaultConfigFilePath string = ""

// ServerConfig struct with fields, useful for configuring metrics server.
type ServerConfig struct {
	StoreInterval  time.Duration `env:"STORE_INTERVAL" json:"store_interval"` // when to flush metrics to disk.
	Address        string        `env:"ADDRESS" json:"address"`               // server address to bind to.
	StoreFile      string        `env:"STORE_FILE" json:"store_file"`         // path to file where metrics are stored.
	Restore        bool          `env:"RESTORE" json:"restore"`               // If set to true, read StoreFile to restore metrics state
	HashKey        string        `env:"KEY" json:"crypto_key"`                // key to create/validate hash
	DBDsn          string        `env:"DATABASE_DSN" json:"database_dsn"`     // Postgres connection string in DSN format
	configFilePath string
}

// NewServerConfig constructor returns pointer to new ServerConfig
func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (c *ServerConfig) mergeConfigs(s ServerConfig) {
	if c.Restore == serverDefaultRestoreCondition {
		c.Restore = s.Restore
	}
	if c.DBDsn == serverDefaultDBDsn {
		c.DBDsn = s.DBDsn
	}
	if c.HashKey == serverDefaultHashKey {
		c.HashKey = s.HashKey
	}
	if c.StoreFile == serverDefaultStoreFile {
		c.StoreFile = s.StoreFile
	}
	if c.StoreInterval == serverDefaultStoreInterval {
		c.StoreInterval = s.StoreInterval
	}
	if c.Address == serverDefaultAddress {
		c.Address = s.Address
	}
}

func (c *ServerConfig) loadConfigFromFile(filename string) (*ServerConfig, error) {
	data, err := os.ReadFile(filename)
	cfg := &ServerConfig{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Parse method to fulfill ServerConfig fields.
// It reads flags and env variables.
func (c *ServerConfig) Parse() error {

	flag.DurationVar(&c.StoreInterval, "i", serverDefaultStoreInterval, "when to flush metrics to disk. Inactive for agent.")
	flag.StringVar(&c.Address, "a", serverDefaultAddress, "http address in format localhost:8080")
	flag.StringVar(&c.StoreFile, "f", serverDefaultStoreFile, "path to file where metrics are stored. Inactive for agent.")
	flag.BoolVar(&c.Restore, "r", serverDefaultRestoreCondition, "If set to true, read file in -f flag to restore metrics state")
	flag.StringVar(&c.HashKey, "k", serverDefaultHashKey, "key to create/validate hash")
	flag.StringVar(&c.DBDsn, "d", serverDefaultDBDsn, "Postgres connection string")
	flag.StringVar(&c.configFilePath, "c", serverDefaultConfigFilePath, "Config file path")
	flag.Parse()

	if c.configFilePath != "" {
		cfg, err := c.loadConfigFromFile(c.configFilePath)
		if err != nil {
			log.Panic("cannot parse config file ", err)
		}
		c.mergeConfigs(*cfg)
	}

	err := env.Parse(c)
	if err != nil {
		return err
	}
	if c.DBDsn != "" {
		c.Restore = false
		c.StoreFile = ""
	}
	return nil
}
