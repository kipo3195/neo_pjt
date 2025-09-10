package repository

import (
	"common/internal/domain/worksInfo/entity"
	"common/internal/domain/worksInfo/models"
	"common/internal/domain/worksInfo/repository"
	"log"

	"gorm.io/gorm"
)

type worksInfoRepositoryImpl struct {
	db *gorm.DB
}

func NewWorksInfoRepository(db *gorm.DB) repository.WorksInfoRepository {
	return &worksInfoRepositoryImpl{
		db: db,
	}
}

func WorksInfoMigrate(db *gorm.DB) {
	db.AutoMigrate(&models.ConnectInfo{})
}

func (r *worksInfoRepositoryImpl) GetConnectInfo(worksCode string) (*entity.ConnectInfo, error) {

	// model
	var connectInfo models.ConnectInfo

	// domain으로 auth에 접근할 것인가?
	result := r.db.Where("works_code = ?", worksCode).First(&connectInfo)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetConnectInfo] - DB error")
		return nil, result.Error
	} else {

		if result.RowsAffected > 0 {
			return &entity.ConnectInfo{
				ServerUrl: connectInfo.ServerUrl,
			}, nil
		} else {
			log.Println("[GetConnectInfo] - DB select X")
			return nil, nil
		}
	}

}
