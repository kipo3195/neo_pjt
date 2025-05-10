package usecases

import (
	"core/dto"
	"core/entities"
	"core/repositories"
	"encoding/json"
	"net/http"
)

type coreUsecase struct {
	repo repositories.CoreRepository
}

type CoreUsecase interface {
	GetValidationData(r *http.Request) (dto.ValidationRequestHeader, dto.ValidationRequest, error)

	CheckValidation(header dto.ValidationRequestHeader) bool
	ToValidationWhereEntity(header dto.ValidationRequestHeader) entities.ValidationWhere

	GetWorksInfo(body dto.ValidationRequest) (entities.WorksInfo, error)
}

func NewCoreUsecase(repo repositories.CoreRepository) CoreUsecase {
	return &coreUsecase{repo: repo}
}

func (u *coreUsecase) GetValidationData(r *http.Request) (dto.ValidationRequestHeader, dto.ValidationRequest, error) {

	validationHeader := dto.ValidationRequestHeader{
		Hash:     r.Header.Get("Hash"),
		Device:   r.Header.Get("Device"),
		Version:  r.Header.Get("Version"),
		Works:    r.Header.Get("Works"),
		DeviceId: r.Header.Get("Device_Id"),
	}

	var validationRequest dto.ValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&validationRequest); err != nil {
		return dto.ValidationRequestHeader{}, dto.ValidationRequest{}, err
	}

	return validationHeader, validationRequest, nil
}

func (u *coreUsecase) CheckValidation(header dto.ValidationRequestHeader) bool {

	validationWhere := u.ToValidationWhereEntity(header)
	flag, err := u.repo.GetValidation(validationWhere)
	if !flag || err != nil {
		return false
	}
	return true
}

func (u *coreUsecase) ToValidationWhereEntity(header dto.ValidationRequestHeader) entities.ValidationWhere {
	return entities.ValidationWhere{
		Hash: header.Hash,
	}
}

func (u *coreUsecase) GetWorksInfo(body dto.ValidationRequest) (entities.WorksInfo, error) {
	result, err := u.repo.GetWorksInfo(body.Domain)
	if err != nil {
		return entities.WorksInfo{}, err
	} else {
		return result, nil
	}

}
