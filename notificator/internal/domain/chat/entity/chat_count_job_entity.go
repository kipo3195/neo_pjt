package entity

import "time"

type ChatCountJobEntity struct {
	UserHash string
	Timer    *time.Timer
	Count    int
	Delta    int
}
