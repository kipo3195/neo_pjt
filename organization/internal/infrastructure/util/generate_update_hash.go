package utils

import (
	"org/internal/infrastructure/consts"
	"time"
)

func GenerateUpdateHash() string {
	return time.Now().Format(consts.YYYYMMDDHHMSS)
}
