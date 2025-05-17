package repositories

import (
	coreerrors "core/consts"
	"core/dto"
	"core/entities"
	"core/models"
	"log"

	"gorm.io/gorm"
)

const (
	DOMAIN = "domain"
	CODE   = "code"
)

type coreRepository struct {
	db *gorm.DB
}

type CoreRepository interface {
	GetValidation(where entities.ValidationWhere) (bool, error)
	GetWorksInfo(body dto.AppValidationRequest) (*entities.WorksInfo, error)
}

func NewCoreRepository(db *gorm.DB) CoreRepository {

	return &coreRepository{db: db}
}

func (r *coreRepository) GetValidation(where entities.ValidationWhere) (bool, error) {
	var validation models.AppValidation

	result := r.db.Where("app_hash = ?", where.Hash).First(&validation)

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

func (r *coreRepository) GetWorksInfo(data dto.AppValidationRequest) (*entities.WorksInfo, error) {

	// model
	var worksInfo models.WorksInfo

	var sql = ""
	var param interface{}

	if data.Type == DOMAIN {
		sql = "works_domain = ? and use_yn = ?"
		param = data.Domain
	} else if data.Type == CODE {
		sql = "works_code = ? and use_yn = ?"
		param = data.Code
	} else {
		log.Println("[GetWorksInfo] - repo type invalid")
		return &entities.WorksInfo{}, coreerrors.ErrInvalidType
	}

	result := r.db.Where(sql, param, "Y").First(&worksInfo)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetWorksInfo] - DB error")
		return &entities.WorksInfo{}, result.Error
	}

	if result.RowsAffected > 0 {
		return &entities.WorksInfo{
			WorksCode: worksInfo.Code,
			WorksName: worksInfo.Name,
			UseYn:     worksInfo.UseYn,
		}, nil
	} else {
		log.Println("[GetWorksInfo] - DB select X")
		return &entities.WorksInfo{}, nil
	}

}
