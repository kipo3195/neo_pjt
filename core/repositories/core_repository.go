package repositories

import (
	coreerrors "core/consts"
	"core/dto"
	"core/entities"
	"core/models"
	"errors"
	"fmt"
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

	result := r.db.Where("app_hash = ?", where.Hash).Where("device_kind = ?", where.Device).First(&validation)

	fmt.Println("appvalidation 에러 여부 ", result.Error)

	// 에러 처리
	if result.Error != nil {
		return false, result.Error
	}
	fmt.Println("appvalidation 조회 수  ", result.RowsAffected)
	if result.RowsAffected > 0 {
		// 1개이상 조회시만 true
		return true, nil
	} else {
		return false, nil
	}

}

func (r *coreRepository) GetWorksInfo(data dto.AppValidationRequest) (*entities.WorksInfo, error) {

	// model
	var worksList models.WorksList

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

	result := r.db.Where(sql, param, "Y").First(&worksList)

	fmt.Println("도메인이나 코드를 전달 받아서 등록된 테넌트 인지 조회 결과 : ", result)

	// 조회 결과가 없을때
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetWorksInfo] - No record found")
		return &entities.WorksInfo{}, coreerrors.ErrInvalidMappingServer
	}

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetWorksInfo] - DB error")
		return &entities.WorksInfo{}, result.Error
	}

	if result.RowsAffected > 0 {
		return &entities.WorksInfo{
			ConnectInfo: entities.ConnectInfo{
				ServerUrl: worksList.ServerUrl,
			},
			WorksCode: worksList.Code,
			WorksName: worksList.Name,
			UseYn:     worksList.UseYn,
		}, nil
	} else {
		log.Println("[GetWorksInfo] - DB select X")
		return &entities.WorksInfo{}, nil
	}

}
