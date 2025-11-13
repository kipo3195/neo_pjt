package usecase

import (
	"notificator/internal/domain/note/repository"
)

type noteUsecase struct {
	repo repository.NoteRepository
}

type NoteUsecase interface {
}

func NewNoteUsecase(repo repository.NoteRepository) NoteUsecase {
	return &noteUsecase{
		repo: repo,
	}
}
