package chatLineService

type ChatLIneEventDto struct {
	EventType     string `json:"eventType"`
	Cmd           int    `json:"cmd"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey,omitempty"`
	Contents      string `json:"contents"`
	SendUserHash  string `json:"sendUserHash"`
	SendDate      string `json:"sendDate"`
	FileId        string `json:"fileId,omitempty"`
	FileName      string `json:"fileName,omitempty"`
	FileType      string `json:"fileType,omitempty"`
	ThumbnailUrl  string `json:"thumbnailUrl,omitempty"`
}
