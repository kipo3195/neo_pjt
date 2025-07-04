package repositories

import (
	coreerrors "core/consts"
	clDto "core/dto/client"
	dto "core/dto/client"
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
	GetWorksCommonInfo(body clDto.AppValidationRequest) (*entities.WorksCommonInfo, error)
}

func NewCoreRepository(db *gorm.DB) CoreRepository {

	return &coreRepository{db: db}
}

func (r *coreRepository) GetValidation(where entities.ValidationWhere) (bool, error) {
	var validation models.AppValidation

	result := r.db.Where("app_hash = ?", where.Hash).Where("device_kind = ?", where.Device).First(&validation)

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetValidation] - No record found or DB error")
		return false, result.Error
	}
	fmt.Println("appvalidation 조회 수  ", result.RowsAffected)

	if result.RowsAffected > 0 {
		// 1개이상 조회시만 true
		return true, nil
	} else {
		log.Println("[GetValidation] - DB select X")
		return false, nil
	}

}

func (r *coreRepository) GetWorksCommonInfo(data dto.AppValidationRequest) (*entities.WorksCommonInfo, error) {

	// model
	var worksList models.WorksList

	result := r.db.Where("works_code = ? and use_yn = ?", data.WorksCode, "Y").First(&worksList)

	fmt.Println("도메인이나 코드를 전달 받아서 등록된 테넌트 인지 조회 결과 : ", result)

	// 조회 결과가 없을때
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetWorksInfo] - No record found")
		return nil, coreerrors.ErrInvalidMappingServer
	}

	// 에러 처리
	if result.Error != nil {
		log.Println("[GetWorksInfo] - DB error")
		return nil, result.Error
	}

	if result.RowsAffected > 0 {
		return &entities.WorksCommonInfo{
			ServerUrl: worksList.ServerUrl,
			WorksCode: worksList.Code,
			WorksName: worksList.Name,
			UseYn:     worksList.UseYn,
		}, nil
	} else {
		log.Println("[GetWorksInfo] - DB select X")
		return nil, nil
	}

}
