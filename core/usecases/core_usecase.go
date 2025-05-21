package usecases

import (
	"bytes"
	consts "core/consts"
	"core/dto"
	"core/entities"
	"core/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type coreUsecase struct {
	repo repositories.CoreRepository
}

type CoreUsecase interface {
	CheckValidation(header *dto.AppValidationRequestHeader) bool

	GetWorksInfo(body dto.AppValidationRequest, uuid string) (*entities.WorksInfo, *dto.ErrorResponse)

	GetConnectInfo(uuid string, serverDomain string) (string, error)

	// 변환
	ToValidationWhereEntity(header *dto.AppValidationRequestHeader) entities.ValidationWhere
}

func NewCoreUsecase(repo repositories.CoreRepository) CoreUsecase {
	return &coreUsecase{repo: repo}
}

func (u *coreUsecase) CheckValidation(header *dto.AppValidationRequestHeader) bool {

	validationWhere := u.ToValidationWhereEntity(header)

	flag, err := u.repo.GetValidation(validationWhere)

	if !flag || err != nil {
		return false
	}

	return true
}

func (u *coreUsecase) ToValidationWhereEntity(header *dto.AppValidationRequestHeader) entities.ValidationWhere {
	return entities.ValidationWhere{
		Hash: header.Hash,
	}
}

func (u *coreUsecase) GetWorksInfo(body dto.AppValidationRequest, uuid string) (*entities.WorksInfo, *dto.ErrorResponse) {

	// 에러타입이 뭐냐에따라 처리..
	result, err := u.repo.GetWorksInfo(body)

	if err != nil {

		switch {
		case errors.Is(err, consts.ErrInvalidType):
			// 타입에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_101,
				Message: consts.E_101_MSG,
			}
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}
		}

	} else {

		// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.

		worksAuth := entities.WorksAuth{}
		connectInfo, err := u.GetConnectInfo(uuid, result.ConnectInfo.ServerDomain)

		if err != nil {
			fmt.Println("에러")
		}

		result.WorksAuth = worksAuth
		result.ConnectInfo.ServerDomain = connectInfo

		return result, nil
	}

}

func (u *coreUsecase) GetConnectInfo(uuid string, serverDomain string) (string, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// POST 요청 보내기
	url := serverDomain + "/common/v1/device-init" // http://localhost:8086
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// 응답 출력
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	fmt.Println("Response:", result)
	return result["result"].(string), nil
}
