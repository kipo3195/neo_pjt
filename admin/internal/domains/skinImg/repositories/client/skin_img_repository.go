package client

import "gorm.io/gorm"

type skingImgRepository struct {
	db *gorm.DB
}

type SkingImgRepository interface {
}

func NewSkinImgRepository(db *gorm.DB) SkingImgRepository {

	return &skingImgRepository{db: db}

}
