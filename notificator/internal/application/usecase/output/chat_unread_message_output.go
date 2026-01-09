package output

type ChatUnreadMessageOutput struct {
	Type           string               `json:"type"`
	EventType      string               `json:"eventType"`
	ChatUnreadData ChatUnreadDataOutput `json:"chatUnreadData"`
}
