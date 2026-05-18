package chatService

type BulkSendChatResponse struct {
	TotalCount int `json:"totalCount"`
	PerSecond  int `json:"perSecond"`
}
