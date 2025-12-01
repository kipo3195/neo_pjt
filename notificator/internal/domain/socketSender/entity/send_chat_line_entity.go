package entity

type SendChatLineEntity struct {
	Cmd          int
	Contents     string
	LineKey      string
	SendUserHash string
	SendDate     string
}

func MakeSendChatLineEntity(cmd int, contents string, lineKey string, sendUserHash string, sendDate string) SendChatLineEntity {
	return SendChatLineEntity{
		Cmd:          cmd,
		Contents:     contents,
		LineKey:      lineKey,
		SendUserHash: sendUserHash,
		SendDate:     sendDate,
	}
}
