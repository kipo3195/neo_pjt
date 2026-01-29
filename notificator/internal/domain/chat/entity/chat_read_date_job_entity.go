package entity

import "time"

type ChatReadDateJobEntity struct {
	UserHash    string
	Timer       *time.Timer
	Count       int64
	RoomReadMap map[string]map[string]string
}
