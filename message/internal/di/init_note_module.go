package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"

	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type NoteModule struct {
	Handler *handler.NoteHandler
	Usecase usecase.NoteUsecase
}

func InitNoteModule(db *gorm.DB, connector *nats.Conn) *NoteModule {

	repository := repository.NewNoteRepository(db)
	usecase := usecase.NewNoteUsecase(repository, connector)
	handler := handler.NewNoteHandler(usecase)

	return &NoteModule{
		Handler: handler,
		Usecase: usecase,
	}
}
