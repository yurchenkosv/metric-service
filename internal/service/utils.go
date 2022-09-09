package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func signHash(key string, msg string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}
