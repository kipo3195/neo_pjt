package input

type SendChatInput struct {
	SendUserHash string
	LineKey      string
	Contents     string
	DestUsers    []string
}
