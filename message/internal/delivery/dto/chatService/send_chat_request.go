package chatService

type SendChatRequest struct {
	Contents  string   `json:"contents"`
	DestUsers []string `json:"destUsers"`
}
