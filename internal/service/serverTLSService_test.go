package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yurchenkosv/metric-service/internal/config"
	"testing"
	"crypto/rsa"
	"crypto/rand"
	log "github.com/sirupsen/logrus"
	"math/big"
)

func getPrivateKey() *rsa.PrivateKey {
	privKey, err := rsa.GenerateKey(rand.Reader,2048)
	if err != nil {
		log.Error(err)
	}
	return privKey
}


func TestNewServerTLSService(t *testing.T) {
	type args struct {
		serverConfig config.ServerConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *ServerTLSService
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewServerTLSService(tt.args.serverConfig)
			if !tt.wantErr(t, err, fmt.Sprintf("NewServerTLSService(%v)", tt.args.serverConfig)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewServerTLSService(%v)", tt.args.serverConfig)
		})
	}
}

func TestServerTLSService_CreatePemCertificateFromPrivateKey(t *testing.T) {
	type fields struct {
		privateKey *rsa.PrivateKey
		cert       []byte
		cfg        config.ServerConfig
	}
	type args struct {
		dnsName []string
		serial *big.Int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		//{
		//	name:    "should successfuly create certificate in pem format",
		//	fields:  fields{
		//		privateKey: getPrivateKey(),
		//	},
		//	args:    args{
		//		dnsName: []string{"tst.server.com"},
		//		serial:  big.NewInt(999),
		//	},
		//	want:    []byte{},
		//	wantErr: assert.NoError,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerTLSService{
				privateKey: tt.fields.privateKey,
				cert:       tt.fields.cert,
				cfg:        tt.fields.cfg,
			}
			got, err := s.CreatePemCertificateFromPrivateKey(tt.args.dnsName...)
			if !tt.wantErr(t, err, fmt.Sprintf("CreatePemCertificateFromPrivateKey(%v)", tt.args.dnsName)) {
				return
			}
			assert.Equalf(t, tt.want, got, "CreatePemCertificateFromPrivateKey(%v)", tt.args.dnsName)
		})
	}
}

func TestServerTLSService_GetPrivateKeyPem(t *testing.T) {
	type fields struct {
		privateKey *rsa.PrivateKey
		cert       []byte
		cfg        config.ServerConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerTLSService{
				privateKey: tt.fields.privateKey,
				cert:       tt.fields.cert,
				cfg:        tt.fields.cfg,
			}
			assert.Equalf(t, tt.want, s.GetPrivateKeyPem(), "GetPrivateKeyPem()")
		})
	}
}

func TestServerTLSService_SaveCertificateToDisk(t *testing.T) {
	type fields struct {
		privateKey *rsa.PrivateKey
		cert       []byte
		cfg        config.ServerConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ServerTLSService{
				privateKey: tt.fields.privateKey,
				cert:       tt.fields.cert,
				cfg:        tt.fields.cfg,
			}
			got, err := s.SaveCertificateToDisk()
			if !tt.wantErr(t, err, fmt.Sprintf("SaveCertificateToDisk()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "SaveCertificateToDisk()")
		})
	}
}

func Test_loadPrivateKey(t *testing.T) {
	type args struct {
		keyPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *rsa.PrivateKey
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadPrivateKey(tt.args.keyPath)
			if !tt.wantErr(t, err, fmt.Sprintf("loadPrivateKey(%v)", tt.args.keyPath)) {
				return
			}
			assert.Equalf(t, tt.want, got, "loadPrivateKey(%v)", tt.args.keyPath)
		})
	}
}
