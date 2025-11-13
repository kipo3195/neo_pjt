package input

type ChatMessageInput struct {
	Type         string
	SendUserHash string
	Contents     string
	LineKey      string
	DestUserHash []string
}
