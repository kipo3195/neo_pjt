package output

type ChatLineDataOutput struct {
	Cmd           int
	Contents      string
	LineKey       string
	TargetLineKey string
	SendUserHash  string
	SendDate      string
}
