package repositories

import "gorm.io/gorm"

type noteRepository struct {
	db *gorm.DB
}

type NoteRepository interface {
}

func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db: db}
}
