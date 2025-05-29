package repositories

import (
	"gorm.io/gorm"
)

type adminRepository struct {
	db *gorm.DB
}

type AdminRepository interface {
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}
