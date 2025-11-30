package entity

type ChatLineEntity struct {
	Cmd          int
	Contents     string
	LineKey      string
	SendUserHash string
	SendDate     string
}

func MakeChatLineEntity(cmd int, contents string, lineKey string, sendUserHash string, sendDate string) ChatLineEntity {
	return ChatLineEntity{
		Cmd:          cmd,
		Contents:     contents,
		LineKey:      lineKey,
		SendUserHash: sendUserHash,
		SendDate:     sendDate,
	}

}
