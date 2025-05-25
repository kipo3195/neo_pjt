package repositories

import (
	"gorm.io/gorm"
)

const (
	DOMAIN = "domain"
	CODE   = "code"
)

type orgRepository struct {
	db *gorm.DB
}

type OrgRepository interface {
}

func NewOrgRepository(db *gorm.DB) OrgRepository {
	return &orgRepository{db: db}
}
