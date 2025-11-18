package usecase

import (
	"context"
	"encoding/json"
	"log"
	"message/internal/application/usecase/input"
	"message/internal/consts"
	"message/internal/domain/note/entity"
	"message/internal/domain/note/repository"

	"github.com/nats-io/nats.go"
)

type noteUsecase struct {
	repository repository.NoteRepository
	connector  *nats.Conn
}

type NoteUsecase interface {
	SendNote(context context.Context, input input.SendNoteInput) error
}

func NewNoteUsecase(repository repository.NoteRepository, connector *nats.Conn) NoteUsecase {
	return &noteUsecase{
		repository: repository,
		connector:  connector,
	}
}

func (u *noteUsecase) SendNote(context context.Context, input input.SendNoteInput) error {

	entity := entity.MakeSendNoteEntity(input.Type, input.NoteKey, input.Contents, input.SendUserHash, input.RecvUserHash, input.RefeUserHash)

	data, err := json.Marshal(entity) // 🔹 struct → []byte(JSON)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// 쪽지 발송
	err = u.connector.Publish("note.message", data)
	if err != nil {
		log.Fatal("NATS publish failed:", err)
		return consts.ErrPublishToMessageBrokerError
	}

	return nil
}
