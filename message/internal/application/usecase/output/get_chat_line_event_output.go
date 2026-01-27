package output

type GetChatLineEventOutput struct {
	EventType     string
	Cmd           int
	LineKey       string
	TargetLineKey string
	Contents      string
	SendUserHash  string
	SendDate      string
	FileId        string
	FileName      string
	FileType      string
	ThumbnailUrl  string
}
