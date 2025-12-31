package dto

type SendChatLineDto struct {
	Cmd           int    `json:"cmd"`
	Contents      string `json:"contents"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	SendUserHash  string `json:"sendUserHash"`
	SendDate      string `json:"sendDate"`
}

func MakeChatLineDto(cmd int, contents string, lineKey string, targetLineKey string, sendUserHash string, sendDate string) SendChatLineDto {
	return SendChatLineDto{
		Cmd:           cmd,
		Contents:      contents,
		LineKey:       lineKey,
		TargetLineKey: targetLineKey,
		SendUserHash:  sendUserHash,
		SendDate:      sendDate,
	}
}
