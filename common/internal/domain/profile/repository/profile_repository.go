package repository

import "gorm.io/gorm"

type profileRepositroy struct {
	db *gorm.DB
}

type ProfileRepository interface {
}
