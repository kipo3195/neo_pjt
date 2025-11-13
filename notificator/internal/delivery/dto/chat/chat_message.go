package chat

type ChatMessage struct {
	Type         string   `json:"type"`
	SendUserHash string   `json:"sendUserHash"`
	Contents     string   `json:"contents"`
	DestUserHash []string `json:"destUserHash"`
}
