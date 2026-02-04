package mapper

import (
	"message/internal/application/usecase/input"
)

func MakeSendNoteInput(sendUserHash string, noteKey string, contents string, t string, recvUserHash []string, refeUserHash []string) input.SendNoteInput {
	return input.SendNoteInput{
		SendUserHash: sendUserHash,
		NoteKey:      noteKey,
		Contents:     contents,
		Type:         t,
		RecvUserHash: recvUserHash,
		RefeUserHash: refeUserHash,
	}
}
