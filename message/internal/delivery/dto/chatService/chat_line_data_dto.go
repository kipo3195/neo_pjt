package chatService

type ChatLineData struct {
	Contents      string `json:"contents"`
	Cmd           int    `json:"cmd"`
	TargetLineKey string `json:"targetLineKey"`
	LineKey       string `json:"lineKey"`      // 서버 생성
	SendUserHash  string `json:"sendUserHash"` // 서버 생성
	SendDate      string `json:"sendDate"`     // 서버 생성
}
