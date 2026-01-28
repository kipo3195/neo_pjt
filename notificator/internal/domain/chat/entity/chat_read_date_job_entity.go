package entity

import "time"

type ChatReadDateJobEntity struct {
	UserHash    string
	Timer       *time.Timer
	RoomReadMap map[string][]ChatReadDateEntity
}
