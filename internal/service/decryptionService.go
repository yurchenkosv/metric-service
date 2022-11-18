package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type DecryptionService struct {
	privateKey *rsa.PrivateKey
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

func NewDecryptionService(keyPath string) (*DecryptionService, error) {
	cert, err := loadPrivateKey(keyPath)
	if err != nil {
		return nil, err
	}
	return &DecryptionService{
			privateKey: cert,
		},
		nil
}

func (s *DecryptionService) DecryptMessage(msg string) (string, error) {
	plainMessage, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		s.privateKey,
		pemStringToCipher(msg),
		nil,
	)
	if err != nil {
		return "", nil
	}
	return string(plainMessage), nil
}

func pemStringToCipher(msg string) []byte {
	b, _ := pem.Decode([]byte(msg))
	return b.Bytes
}
