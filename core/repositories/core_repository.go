package repositories

import (
	"core/entities"
	"core/models"

	"gorm.io/gorm"
)

type coreRepository struct {
	db *gorm.DB
}

type CoreRepository interface {
	GetValidation(where entities.ValidationWhere) (bool, error)
	GetWorksInfo(body string) (entities.WorksInfo, error)
}

func NewCoreRepository(db *gorm.DB) CoreRepository {

	return &coreRepository{db: db}
}

func (r *coreRepository) GetValidation(where entities.ValidationWhere) (bool, error) {
	var validation models.AppValidation

	result := r.db.Where("version_id = ?", where.Hash).First(&validation)

	// 에러 처리
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		// 1개이상 조회시만 true
		return true, nil
	} else {
		return false, nil
	}

}

func (r *coreRepository) GetWorksInfo(domain string) (entities.WorksInfo, error) {

	var worksInfo models.WorksInfo

	result := r.db.Where("w_domain = ? and use_yn = ?", domain, "Y").First(&worksInfo)

	// model -> entity 변환처리

	// 에러 처리
	if result.Error != nil {
		return entities.WorksInfo{}, result.Error
	}

	if result.RowsAffected > 0 {
		return entities.WorksInfo{
			WCode:   worksInfo.WCode,
			WName:   worksInfo.WName,
			RegDate: worksInfo.RegDate,
			WDomain: worksInfo.WDomain,
		}, nil
	} else {
		return entities.WorksInfo{}, nil
	}

}
