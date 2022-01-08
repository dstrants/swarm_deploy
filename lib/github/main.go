package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Calculates the HMAC from a given data array and secret string
func calculateHMAC(data []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
