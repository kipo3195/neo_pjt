package di

import (
	"notificator/internal/infrastructure/broker"
	"notificator/internal/infrastructure/repository"

	"notificator/internal/application/usecase"

	"gorm.io/gorm"
)

type NoteModule struct {
	Usecase usecase.NoteUsecase
}

func InitNoteModule(db *gorm.DB, mb broker.Broker) *NoteModule {

	repo := repository.NewNoteRepository(db)
	usecase := usecase.NewNoteUsecase(repo, mb)

	return &NoteModule{
		Usecase: usecase,
	}
}
