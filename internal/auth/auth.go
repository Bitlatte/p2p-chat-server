package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateUserID(value1, value2 string) string {
	data := value1 + value2
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
