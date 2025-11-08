package usecases

import (
	"message/repositories"

	"github.com/gorilla/websocket"
)

type noteUsecase struct {
	repo repositories.NoteRepository
}

type NoteUsecase interface {
	HandleNote(conn *websocket.Conn, data map[string]interface{})
}

func NewNoteUsecase(repo repositories.NoteRepository) NoteUsecase {
	return &noteUsecase{repo: repo}
}

func (r *noteUsecase) HandleNote(conn *websocket.Conn, data map[string]interface{}) {

}
