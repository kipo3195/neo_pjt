package input

type BulkSendChatInput struct {
	TotalCount  int
	PerSecond   int
	WorkerCount int
	Contents    string
	Cmd         int
	EventType   string
	ChatSession string
}
