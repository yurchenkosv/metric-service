package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

type EncryptionService struct {
	publicKey *rsa.PublicKey
}

func loadPublicKey(keyPath string) (*rsa.PublicKey, error) {
	bytes, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(bytes)
	blockBytes := block.Bytes

	publicKey, err := x509.ParsePKCS1PublicKey(blockBytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func NewEncryptionService(keyPath string) (*EncryptionService, error) {
	pubKey, err := loadPublicKey(keyPath)
	if err != nil {
		return nil, err
	}
	return &EncryptionService{
			publicKey: pubKey,
		},
		nil
}

func (s EncryptionService) EncryptMessage(msg string) (string, error) {
	cipher, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		s.publicKey,
		[]byte(msg),
		nil)
	if err != nil {
		return "", err
	}

	return cipherToPemString(cipher), nil
}

func cipherToPemString(cipher []byte) string {
	return string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "MESSAGE",
				Bytes: cipher,
			},
		),
	)
}
