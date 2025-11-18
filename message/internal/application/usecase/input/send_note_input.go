package input

type SendNoteInput struct {
	SendUserHash string
	NoteKey      string
	Contents     string
	Type         string
	RecvUserHash []string
	RefeUserHash []string
}
