package entity

type SendNoteEntity struct {
	Type         string   `json:"type"`
	SendUserHash string   `json:"sendUserHash"`
	Contents     string   `json:"contents"`
	NoteKey      string   `json:"noteKey"`
	RecvUserHash []string `json:"recvUserHash"`
	RefeUserHash []string `json:"refeUserHash"`
}

func MakeSendNoteEntity(t string, noteKey string, contents string, sendUserHash string, recvUserHash []string, refeUserHash []string) SendNoteEntity {

	return SendNoteEntity{
		Type:         t,
		SendUserHash: sendUserHash,
		Contents:     contents,
		NoteKey:      noteKey,
		RecvUserHash: recvUserHash,
		RefeUserHash: refeUserHash,
	}
}
