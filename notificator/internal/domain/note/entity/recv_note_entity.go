package entity

type RecvNoteMessage struct {
	Type         string   `json:"type"`
	NoteKey      string   `json:"noteKey"`
	Contents     string   `json:"contents"`
	SendUserHash string   `json:"sendUserHash"`
	RecvUserHash []string `json:"recvUserHash"`
	RefeUserHash []string `json:"refeUserHash"`
}

func MakeRecvNoteMessageEntity(t string, noteKey string, contents string, sendUserHash string, recv []string, refe []string) RecvNoteMessage {

	return RecvNoteMessage{
		Type:         t,
		NoteKey:      noteKey,
		Contents:     contents,
		SendUserHash: sendUserHash,
		RecvUserHash: recv,
		RefeUserHash: refe,
	}

}
