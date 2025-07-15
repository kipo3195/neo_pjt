package repositories

import (
	"core/entities"
	"core/models"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type coreRepository struct {
	db *gorm.DB
}

type CoreRepository interface {
	GetValidation(where entities.ValidationEntity) (bool, error)
	GetWorksCommonInfo(worksCode string) (*entities.WorksCommonInfo, error)
}

func NewCoreRepository(db *gorm.DB) CoreRepository {

	return &coreRepository{db: db}
}

func (r *coreRepository) GetValidation(where entities.ValidationEntity) (bool, error) {
	var validation models.AppValidation

	result := r.db.Where("app_hash = ?", where.Hash).Where("device_kind = ?", where.Device).First(&validation)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 조회 결과가 없을 때
		fmt.Print("[GetValidation] result record = 0")
		return false, result.Error
	} else if result.Error != nil {
		// 기타 에러 발생 (DB 오류 등)
		fmt.Print("[GetValidation] DB error")
		return false, result.Error
	}

	return true, nil
}

func (r *coreRepository) GetWorksCommonInfo(worksCode string) (*entities.WorksCommonInfo, error) {

	// model
	var worksList models.WorksList

	result := r.db.Where("works_code = ? and use_yn = ?", worksCode, "Y").First(&worksList)

	fmt.Println("도메인이나 코드를 전달 받아서 등록된 테넌트 인지 조회 결과 : ", result)

	// 조회 결과가 없을때
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("[GetWorksCommonInfo] - No record found")
		return nil, result.Error
	} else if result.Error != nil {
		log.Println("[GetWorksCommonInfo] - DB error")
		return nil, result.Error
	}

	return &entities.WorksCommonInfo{
		ServerUrl: worksList.ServerUrl,
		WorksCode: worksList.Code,
		WorksName: worksList.Name,
		UseYn:     worksList.UseYn,
	}, nil

}
