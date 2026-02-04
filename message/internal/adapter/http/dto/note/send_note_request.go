package note

type SendNoteRequest struct {
	Type         string   `json:"type"`
	Contents     string   `json:"contents"`
	NoteKey      string   `json:"noteKey"`
	RecvUserHash []string `json:"recvUserHash"`
	RefeUserHash []string `json:"refeUserHash"`
}
