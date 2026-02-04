package repository

import (
	"message/internal/domain/note/repository"

	"gorm.io/gorm"
)

type noteRepository struct {
	db *gorm.DB
}

func NewNoteRepository(db *gorm.DB) repository.NoteRepository {
	return &noteRepository{
		db: db,
	}
}
