package dto

type ChatCountMessageResponse struct {
	Type          string           `json:"type"`
	EventType     string           `json:"eventType"`
	ChatCountData ChatCountDataDto `json:"chatCountData"`
}
