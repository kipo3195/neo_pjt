package chatService

type SendChatRequest struct {
	Contents string   `json:"contents"`
	DestIds  []string `json:"destIds"`
}
