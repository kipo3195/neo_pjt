package chatService

type BulkSendChatRequest struct {
	TotalCount  int    `json:"totalCount" validate:"required,min=1"`
	PerSecond   int    `json:"perSecond" validate:"required,min=1"`
	WorkerCount int    `json:"workerCount" validate:"omitempty,min=1"`
	Contents    string `json:"contents"`
	Cmd         int    `json:"cmd"`
	EventType   string `json:"eventType"`
	ChatSession string `json:"chatSession"`
}
