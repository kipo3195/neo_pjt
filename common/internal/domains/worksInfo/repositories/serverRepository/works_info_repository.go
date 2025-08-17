package serverRepository

import (
	"common/internal/domains/worksInfo/entities"
	"common/internal/domains/worksInfo/models"
	"log"

	"gorm.io/gorm"
)

type worksInfoRepository struct {
	db *gorm.DB
}

type WorksInfoRepository interface {
	GetConnectInfo(worksCode string) (*entities.ConnectInfo, error)
}

func NewWorksInfoRepository(db *gorm.DB) WorksInfoRepository {
	return &worksInfoRepository{
		db: db,
	}
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.ConnectInfo{})
}

func (r *worksInfoRepository) GetConnectInfo(worksCode string) (*entities.ConnectInfo, error) {

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
			return &entities.ConnectInfo{
				ServerUrl: connectInfo.ServerUrl,
			}, nil
		} else {
			log.Println("[GetConnectInfo] - DB select X")
			return nil, nil
		}
	}

}
