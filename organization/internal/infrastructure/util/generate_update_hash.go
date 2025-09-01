package utils

import (
	"org/internal/consts"
	"time"
)

func GenerateUpdateHash() string {
	return time.Now().Format(consts.YYYYMMDDHHMSS)
}
