package input

type ChatLineDataInput struct {
	Cmd           int    `json:"cmd"`
	SendUserHash  string `json:"sendUserHash"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	Contents      string `json:"contents"`
	SendDate      string `json:"sendDate"`
}
