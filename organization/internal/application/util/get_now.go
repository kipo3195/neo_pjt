package util

import (
	"org/internal/application/consts"
	"time"
)

func GetNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}
