package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/yurchenkosv/metric-service/internal/config"
)

type ServerTLSService struct {
	privateKey *rsa.PrivateKey
	cert       []byte
	cfg        config.ServerConfig
}

func loadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	blockBytes := block.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func NewServerTLSService(serverConfig config.ServerConfig) (*ServerTLSService, error) {
	cert, err := loadPrivateKey(serverConfig.CryptoKey)
	if err != nil {
		return nil, err
	}
	return &ServerTLSService{
			cfg:        serverConfig,
			privateKey: cert,
		},
		nil
}

func (s *ServerTLSService) CreatePemCertificateFromPrivateKey(dnsName ...string) ([]byte, error) {
	maxSerialNum := big.NewInt(9999999)
	rndNum, err := rand.Int(rand.Reader, maxSerialNum)
	if err != nil {
		log.Error("cannot create random serial ", err)
	}
	tml := x509.Certificate{
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(5, 0, 0),
		SerialNumber: rndNum,
		Subject: pkix.Name{
			CommonName:   dnsName[0],
			Organization: []string{"TLS Ltd."},
		},
		DNSNames:              dnsName,
		IsCA:                  false,
		BasicConstraintsValid: true,
	}
	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &s.privateKey.PublicKey, s.privateKey)
	if err != nil {
		return nil, err
	}

	certPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	s.cert = certPem
	return certPem, nil
}

func (s *ServerTLSService) SaveCertificateToDisk() (string, error) {
	storeLocation := filepath.Dir(s.cfg.CryptoKey)
	cryptoFilename := strings.Split(filepath.Base(s.cfg.CryptoKey), ".")[0]
	certLocation := storeLocation + "/" + cryptoFilename + ".crt"
	err := os.WriteFile(certLocation, s.cert, 0600)
	return certLocation, err
}

func (s *ServerTLSService) GetPrivateKeyPem() []byte {
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.privateKey),
	})
	return privateKeyPem
}
