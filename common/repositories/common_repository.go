package repositories

import (
	"common/entities"
	"common/models"
	"log"

	"gorm.io/gorm"
)

type commonRepository struct {
	db *gorm.DB
}

type CommonRepository interface {
	GetConnectInfo(where string) (*entities.InitResult, error)
}

func NewCommonRepository(db *gorm.DB) CommonRepository {

	return &commonRepository{db: db}
}

func (r *commonRepository) GetConnectInfo(worksCode string) (*entities.InitResult, error) {

	// model
	var connectInfo models.ConnectInfo

	// domain으로 auth에 접근할 것인가?
	result := r.db.Where("works_code = ?", worksCode).First(&connectInfo)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetConnectInfo] - DB error")
		return &entities.InitResult{}, result.Error
	} else {

		if result.RowsAffected > 0 {
			return &entities.InitResult{
				ConnectInfo: connectInfo.ServerUrl,
			}, nil
		} else {
			log.Println("[GetConnectInfo] - DB select X")
			return &entities.InitResult{}, nil
		}
	}

}
