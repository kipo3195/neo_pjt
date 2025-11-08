package repository

import "gorm.io/gorm"

type noteRepository struct {
	db *gorm.DB
}

type NoteRepository interface {
}
