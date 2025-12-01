package dto

type ChatLineDto struct {
	Cmd          int    `json:"cmd"`
	Contents     string `json:"contents"`
	LineKey      string `json:"lineKey"`
	SendUserHash string `json:"sendUserHash"`
	SendDate     string `json:"sendDate"`
}

func MakeChatLineDto(cmd int, contents string, lineKey string, sendUserHash string, sendDate string) ChatLineDto {
	return ChatLineDto{
		Cmd:          cmd,
		Contents:     contents,
		LineKey:      lineKey,
		SendUserHash: sendUserHash,
		SendDate:     sendDate,
	}
}
