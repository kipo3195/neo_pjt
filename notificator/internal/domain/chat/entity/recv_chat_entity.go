package entity

type RecvChatMessage struct {
	Type         string `json:"type"`
	LineKey      string `json:"lineKey"`
	Contents     string `json:"contents"`
	SendUserHash string `json:"sendUserHash"`
}

func MakeRecvChatMessageEntity(t string, lineKey string, contents string, sendUserHash string) RecvChatMessage {

	return RecvChatMessage{
		Type:         t,
		LineKey:      lineKey,
		Contents:     contents,
		SendUserHash: sendUserHash,
	}

}
