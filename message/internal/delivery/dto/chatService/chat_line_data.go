package chatService

type ChatLineData struct {
	Contents     string `json:"contents"`
	LineKey      string `json:"lineKey"`
	SendUserHash string `json:"sendUserHash"`
	EventType    string `json:"eventType"`
}
