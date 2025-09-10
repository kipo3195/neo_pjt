package repository

import (
	"common/internal/domain/worksInfo/entity"

	"gorm.io/gorm"
)

type worksInfoRepository struct {
	db *gorm.DB
}

type WorksInfoRepository interface {
	GetConnectInfo(worksCode string) (*entity.ConnectInfo, error)
}
