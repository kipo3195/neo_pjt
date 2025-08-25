package utils

import (
	"org/internal/consts"
	"time"
)

func GetNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}
