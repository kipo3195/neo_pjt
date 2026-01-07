package output

type ChatLineDataOutput struct {
	Cmd           int    `json:"cmd"`
	Contents      string `json:"contents"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	SendUserHash  string `json:"sendUserHash"`
	SendDate      string `json:"sendDate"`
}
