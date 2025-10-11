package repository

import (
	"admin/internal/domain/skinImg/repository"

	"gorm.io/gorm"
)

type skingImgRepositoryImpl struct {
	db *gorm.DB
}

func NewSkinImgRepository(db *gorm.DB) repository.SkinImgRepository {

	return &skingImgRepositoryImpl{db: db}

}
