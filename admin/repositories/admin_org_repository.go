package repositories

import (
	"gorm.io/gorm"
)

type adminOrgRepository struct {
	db *gorm.DB
}

type AdminOrgRepository interface {
}

func NewAdminOrgRepository(db *gorm.DB) AdminOrgRepository {
	return &adminOrgRepository{db: db}
}
