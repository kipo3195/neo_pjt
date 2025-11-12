package chat

type ChatMessage struct {
	Type     string `json:"type"`
	UserHash string `json:"userHash"`
}
