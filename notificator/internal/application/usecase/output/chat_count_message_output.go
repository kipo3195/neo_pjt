package output

type ChatCountMessageOutput struct {
	Type          string              `json:"type"`
	EventType     string              `json:"eventType"`
	ChatCountData ChatCountDataOutput `json:"chatCountData"`
}
