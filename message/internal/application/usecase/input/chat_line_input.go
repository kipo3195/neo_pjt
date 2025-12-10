package input

type ChatLineInput struct {
	Cmd           int
	SendUserHash  string
	LineKey       string
	TargetLineKey string
	Contents      string
	SendDate      string
}
