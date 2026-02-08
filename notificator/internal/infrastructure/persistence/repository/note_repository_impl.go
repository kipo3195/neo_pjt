package repository

import (
	"notificator/internal/domain/note/repository"

	"gorm.io/gorm"
)

type noteRepositoryimpl struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) repository.NoteRepository {
	return &noteRepositoryimpl{db: db}
}
