package entity

type ChatLineEntity struct {
	Cmd           int    `json:"cmd"`
	SendUserHash  string `json:"sendUserHash"`
	LineKey       string `json:"lineKey"`
	TargetLineKey string `json:"targetLineKey"`
	Contents      string `json:"contents"`
	SendDate      string `json:"sendDate"`
}

func MakeChatLineEntity(cmd int, contents string, lineKey string, targetLineKey string, sendUserHash string, sendDate string) ChatLineEntity {
	return ChatLineEntity{
		Cmd:           cmd,
		Contents:      contents,
		LineKey:       lineKey,
		TargetLineKey: targetLineKey,
		SendUserHash:  sendUserHash,
		SendDate:      sendDate,
	}
}
