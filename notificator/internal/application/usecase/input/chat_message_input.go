package input

type ChatMessageInput struct {
	Type         string
	SendUserHash string
	Contents     string
	DestUserHash []string
}
