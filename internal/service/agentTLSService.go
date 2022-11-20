package service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/yurchenkosv/metric-service/internal/config"
)

type AgentTLSService struct {
	cfg  config.AgentConfig
	cert *x509.Certificate
}

func loadCert(keyPath string) (*x509.Certificate, error) {
	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	blockBytes := block.Bytes

	cert, err := x509.ParseCertificate(blockBytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func NewAgentTLSService(cfg config.AgentConfig) (*AgentTLSService, error) {
	cert, err := loadCert(cfg.CryptoKey)
	if err != nil {
		return nil, err
	}
	return &AgentTLSService{
			cert: cert,
		},
		nil
}

func (s *AgentTLSService) GetTLSConfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPool.AddCert(s.cert)
	tlsConfig := tls.Config{RootCAs: certPool}
	return &tlsConfig
}

func (s *AgentTLSService) GetCertificate() *x509.Certificate {
	return s.cert
}
