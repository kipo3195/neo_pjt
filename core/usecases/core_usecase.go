package usecases

import (
	"bytes"
	consts "core/consts"
	clDto "core/dto/client"
	dto "core/dto/common"
	svDto "core/dto/server"
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
	CheckValidation(header *clDto.AppValidationRequestHeader) bool

	GetWorksInfo(body clDto.AppValidationRequest, uuid string) (*entities.WorksInfo, *dto.ErrorResponse, bool)

	GetConnectInfo(uuid string, worksCode string, serverDomain string) (*svDto.SvDeviceInitResponse, error)

	// 변환
	ToValidationWhereEntity(header *clDto.AppValidationRequestHeader) entities.ValidationWhere
}

func NewCoreUsecase(repo repositories.CoreRepository) CoreUsecase {
	return &coreUsecase{repo: repo}
}

func (u *coreUsecase) CheckValidation(header *clDto.AppValidationRequestHeader) bool {

	validationWhere := u.ToValidationWhereEntity(header)

	flag, err := u.repo.GetValidation(validationWhere)

	if !flag || err != nil {
		return false
	}

	return true
}

func (u *coreUsecase) ToValidationWhereEntity(header *clDto.AppValidationRequestHeader) entities.ValidationWhere {
	return entities.ValidationWhere{
		Hash:   header.Hash,
		Device: header.Device,
	}
}

func (u *coreUsecase) GetWorksInfo(body clDto.AppValidationRequest, uuid string) (*entities.WorksInfo, *dto.ErrorResponse, bool) {

	// 에러타입이 뭐냐에따라 처리..
	result, err := u.repo.GetWorksInfo(body)

	fmt.Println("클라이언트가 전달한 도메인, 코드로 전달했을때 조회된 서버 정보 : ", result)

	if err != nil {
		switch {
		case errors.Is(err, consts.ErrInvalidType):
			// 타입에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_101,
				Message: consts.E_101_MSG,
			}, false
		case errors.Is(err, consts.ErrInvalidMappingServer):
			// 매핑된 서버 정보 없음
			return nil, &dto.ErrorResponse{
				Code:    consts.CORE_F102,
				Message: consts.CORE_F102_MSG,
			}, true
		default:
			// 기타 DB 에러
			return nil, &dto.ErrorResponse{
				Code:    consts.E_102,
				Message: consts.E_102_MSG,
			}, false
		}

	} else {

		// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.

		deviceInitResponse, err := u.GetConnectInfo(uuid, body.WorksCode, result.ConnectInfo.ServerUrl)

		if err != nil {
			fmt.Println("common service 호출시 에러 발생함.")
			return &entities.WorksInfo{}, &dto.ErrorResponse{
				Code:    consts.E_500,
				Message: consts.E_500_MSG,
			}, false
		}

		worksAuth := entities.WorksAuth{AppToken: deviceInitResponse.AppToken}
		result.WorksAuth = worksAuth
		result.ConnectInfo.ServerUrl = deviceInitResponse.ServerUrl

		return result, nil, false
	}

}

func (u *coreUsecase) GetConnectInfo(uuid string, worksCode string, serverUrl string) (*svDto.SvDeviceInitResponse, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid":      uuid,
		"worksCode": worksCode,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return &svDto.SvDeviceInitResponse{}, err
	}

	// POST 요청 보내기
	url := "http://" + serverUrl + "/common/v1/device-init"

	fmt.Println("common service 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &svDto.SvDeviceInitResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &svDto.SvDeviceInitResponse{}, err
	}

	defer resp.Body.Close()

	// serverResponse로 전달받기

	var result dto.Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("serverReponse 파싱시 에러")
		return &svDto.SvDeviceInitResponse{}, err
	}

	resultData, ok := result.Data.(map[string]interface{})
	if !ok {
		fmt.Println("Data 필드를 map으로 변환하는 데 실패했습니다.")
		return &svDto.SvDeviceInitResponse{}, errors.New("invalid data format")
	}

	authToken, authTokenOk := resultData["authToken"].(string)
	connectInfo, connectInfoOk := resultData["connectInfo"].(string)

	if !authTokenOk || !connectInfoOk {
		fmt.Println("authToken 또는 connectInfo string으로 변환하는 데 실패했습니다.")
		return &svDto.SvDeviceInitResponse{}, errors.New("invalid token format")
	}

	return &svDto.SvDeviceInitResponse{
		AppToken:  authToken,
		ServerUrl: connectInfo,
	}, nil
}
