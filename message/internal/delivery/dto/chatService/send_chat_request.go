package chatService

type SendChatRequest struct {
	Contents     string   `json:"contents"`
	RecvUserHash []string `json:"recvUserHash"`
}
