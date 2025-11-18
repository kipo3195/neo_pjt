package di

import (
	"notificator/internal/infrastructure/storage"

	"notificator/internal/application/usecase"

	"gorm.io/gorm"
)

type NoteModule struct {
	Usecase usecase.NoteUsecase
}

func InitNoteModule(db *gorm.DB, noteUserStorage storage.NoteUserStorage) *NoteModule {

	//repo := repository.NewNoteRepository(db)
	usecase := usecase.NewNoteUsecase(noteUserStorage)

	return &NoteModule{
		Usecase: usecase,
	}
}
