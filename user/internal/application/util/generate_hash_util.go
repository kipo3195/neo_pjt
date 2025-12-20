package util

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"user/internal/consts"
)

func GenerateUserHashSHA256(userId string) string {
	date := time.Now().Format(consts.YYYYMMDDHHMMSS)
	temp := userId + date
	hash := sha256.Sum256([]byte(temp))
	return hex.EncodeToString(hash[:])
}

func GenerateIntHash() int64 {
	return time.Now().UnixNano()
}
