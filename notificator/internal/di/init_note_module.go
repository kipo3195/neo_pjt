package di

import (
	"notificator/internal/application/usecase"

	"gorm.io/gorm"
)

type NoteModule struct {
	Usecase usecase.NoteUsecase
}

func InitNoteModule(db *gorm.DB) *NoteModule {

	//repo := repository.NewNoteRepository(db)
	usecase := usecase.NewNoteUsecase()

	return &NoteModule{
		Usecase: usecase,
	}
}
