package input

type NoteMessageInput struct {
	Type         string   `json:"type"`
	SendUserHash string   `json:"sendUserHash"`
	Contents     string   `json:"contents"`
	NoteKey      string   `json:"noteKey"`
	RecvUserHash []string `json:"recvUserHash"`
	RefeUserHash []string `json:"refeUserHash"`
}
