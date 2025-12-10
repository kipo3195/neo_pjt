package entity

type ChatLineEntity struct {
	Cmd           int
	Contents      string
	LineKey       string
	TargetLineKey string
	SendUserHash  string
	SendDate      string
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
