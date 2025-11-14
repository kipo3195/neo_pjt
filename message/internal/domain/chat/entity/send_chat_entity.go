package entity

type SendChatEntity struct {
	Type         string   `json:"type"`
	SendUserHash string   `json:"sendUserHash"`
	Contents     string   `json:"contents"`
	LineKey      string   `json:"lineKey"`
	DestUsers    []string `json:"destUsers"`
}

func MakeSendChatEntity(t string, sendUserHash string, contents string, linekey string, destUsers []string) SendChatEntity {

	return SendChatEntity{
		Type:         t,
		SendUserHash: sendUserHash,
		Contents:     contents,
		LineKey:      linekey,
		DestUsers:    destUsers,
	}
}
