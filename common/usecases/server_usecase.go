package usecases

import (
	"bytes"
	"common/consts"
	dto "common/dto/common"
	svDto "common/dto/server"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type serverUsecase struct {
	repo              repositories.ServerRepository
	configHashStorage storage.ConfigHashStorage
}

type ServerUsecase interface {
	DeviceInit(body *svDto.SvDeviceInitRequest, ctx context.Context) (*entities.InitResult, *dto.ErrorResponse)
}

func NewServerUsecase(repo repositories.ServerRepository, configHashStorage storage.ConfigHashStorage) ServerUsecase {
	return &serverUsecase{
		repo:              repo,
		configHashStorage: configHashStorage,
	}
}

func (u *serverUsecase) DeviceInit(body *svDto.SvDeviceInitRequest, ctx context.Context) (*entities.InitResult, *dto.ErrorResponse) {

	// DB 조회
	result, err := u.repo.GetConnectInfo(body.WorksCode)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_102,
			Message: consts.E_102_MSG,
		}
	}

	// AUTH에 JWT 요청
	result.AppToken, err = generateDeviceToken(body, result.ConnectInfo)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// 타임존, 언어, 앱 별 스킨 정보, 설정 정보
	worksConfig, err := u.repo.GetWorksConfig(toWorksConfigEntity(body.WorksCode, body.Device), ctx)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	result.TimeZone = worksConfig.TimeZone
	result.Language = worksConfig.Language
	result.SkinVersion = worksConfig.SkinVersion
	result.ConfigVersion = worksConfig.ConfigVersion

	return result, nil
}

func generateDeviceToken(body *svDto.SvDeviceInitRequest, serverUrl string) (string, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": body.Uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	fmt.Println("auth service 호출! 1")

	url := "http://" + serverUrl + "/auth/v1/generate-device-token"
	//url := domain + "/auth/v1/generate-device-token"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 api key

	// POST 요청 보내기
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("auth service 호출 에러 1")
		return "", err
	}
	defer resp.Body.Close()

	// 구조체로 반환해야 하는거아닌가?
	// 서버간 통신에서 var result dto.ServerResponsed 이 구조를 사용할 것인지 고민

	// 응답 출력
	var result dto.Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("serverReponse 파싱시 에러")
		return "", err
	}

	resultData, ok := result.Data.(map[string]interface{})
	if !ok {
		fmt.Println("Data 필드를 map으로 변환하는 데 실패했습니다.")
		return "", errors.New("invalid data format")
	}

	token, tokenOk := resultData["token"].(string)

	if !tokenOk {
		fmt.Println("token 또는 uuid를 string으로 변환하는 데 실패했습니다.")
		return "", errors.New("invalid token format")
	}
	fmt.Println("auth service 호출 후 발급 받은 토큰 : ", token)
	return token, nil
}
