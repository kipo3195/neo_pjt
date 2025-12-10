package entity

type SendChatLineEntity struct {
	Cmd           int
	Contents      string
	LineKey       string
	TargetLineKey string
	SendUserHash  string
	SendDate      string
}

func MakeSendChatLineEntity(cmd int, contents string, lineKey string, targetLineKey string, sendUserHash string, sendDate string) SendChatLineEntity {
	return SendChatLineEntity{
		Cmd:           cmd,
		Contents:      contents,
		LineKey:       lineKey,
		TargetLineKey: targetLineKey,
		SendUserHash:  sendUserHash,
		SendDate:      sendDate,
	}
}
