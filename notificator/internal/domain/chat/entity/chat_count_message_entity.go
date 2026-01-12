package entity

type ChatCountMessageEntity struct {
	Type          string              `json:"type"`
	EventType     string              `json:"eventType"`
	ChatCountData ChatCountDataEntity `json:"chatCountData"`
}
