package config

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

const agentDefaultAddress string = "localhost:8080"
const agentDefaultReportInterval time.Duration = 10 * time.Second
const agentDefaultPollInterval time.Duration = 2 * time.Second
const agentDefaultHashKey string = ""
const agentDefaultConfigFilepath string = ""
const agentDefaultCryptoKeyPath = ""
const agentDefaultTransportType = "http"

// AgentConfig struct with fields, useful for configuring metrics agent.
type AgentConfig struct {
	Address        string        `env:"ADDRESS"`         // http address of metric server
	ReportInterval time.Duration `env:"REPORT_INTERVAL"` // interval to send metrics to server
	PollInterval   time.Duration `env:"POLL_INTERVAL"`   // interval to collect metrics
	HashKey        string        `env:"KEY"`             // key to create hash
	CryptoKey      string        `env:"CRYPTO_KEY"`      // path to private cipher key
	TransportType  string        `env:"TRANSPORT"`
	configFilePath string
}

func (c *AgentConfig) loadConfigFromFile(filename string) (*AgentConfig, error) {
	data, err := os.ReadFile(filename)
	cfg := &AgentConfig{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *AgentConfig) mergeConfigs(s AgentConfig) {
	if c.Address == agentDefaultAddress {
		c.Address = s.Address
	}
	if c.ReportInterval == agentDefaultReportInterval {
		c.ReportInterval = s.ReportInterval
	}
	if c.PollInterval == agentDefaultPollInterval {
		c.PollInterval = s.PollInterval
	}
	if c.HashKey == agentDefaultHashKey {
		c.HashKey = s.HashKey
	}
	if c.CryptoKey == agentDefaultCryptoKeyPath {
		c.CryptoKey = s.CryptoKey
	}
	if c.TransportType == agentDefaultTransportType {
		c.TransportType = s.TransportType
	}
}

// Parse method to fulfill AgentConfig fields.
// It reads flags and env variables.
func (c *AgentConfig) Parse() error {
	flag.StringVar(&c.Address, "a", agentDefaultAddress, "http address to send metrics in format localhost:8080")
	flag.DurationVar(&c.ReportInterval, "r", agentDefaultReportInterval, "interval to send metrics to server. Inactive for server.")
	flag.DurationVar(&c.PollInterval, "p", agentDefaultPollInterval, "Interval to collect metrics. Inactive for server.")
	flag.StringVar(&c.HashKey, "k", agentDefaultHashKey, "key to create hash")
	flag.StringVar(&c.configFilePath, "c", agentDefaultConfigFilepath, "path to config file")
	flag.StringVar(&c.CryptoKey, "crypto-key", agentDefaultCryptoKeyPath, "path to public key to encrypt messages")
	flag.StringVar(&c.TransportType, "transport", agentDefaultTransportType, "transport type could be http or grpc")
	flag.Parse()

	if c.configFilePath != "" {
		cfg, err := c.loadConfigFromFile(c.configFilePath)
		if err != nil {
			log.Panic("cannot parse config file ", err)
		}
		c.mergeConfigs(*cfg)
	}

	err := env.Parse(c)
	return err
}
