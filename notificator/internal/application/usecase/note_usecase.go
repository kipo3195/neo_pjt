package usecase

import (
	"notificator/internal/domain/note/repository"
	"notificator/internal/infrastructure/broker"
)

type noteUsecase struct {
	repo repository.NoteRepository
	mb   broker.Broker
}

type NoteUsecase interface {
	//HandleNote(conn *websocket.Conn, data map[string]interface{})
}

func NewNoteUsecase(repo repository.NoteRepository, mb broker.Broker) NoteUsecase {
	return &noteUsecase{
		repo: repo,
		mb:   mb,
	}
}

// func (r *noteUsecase) HandleNote(conn *websocket.Conn, data map[string]interface{}) {

// }
