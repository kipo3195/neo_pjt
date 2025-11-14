package chat

type ChatMessage struct {
	Type         string   `json:"type"`
	SendUserHash string   `json:"sendUserHash"`
	Contents     string   `json:"contents"`
	LineKey      string   `json:"lineKey"`
	DestUsers    []string `json:"destUsers"`
}
