package helper

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateTraceId() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
