package usecases

import (
	"bytes"
	consts "common/consts"
	clDto "common/dto/client"
	dto "common/dto/common"
	svDto "common/dto/server"
	"common/entities"
	"common/repositories"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type commonUsecase struct {
	repo repositories.CommonRepository
}

type CommonUsecase interface {
	GetConfig(clDto.ConfigRequest) (*entities.Config, error)
	DeviceInit(body *svDto.SvDeviceInitRequest, ctx context.Context) (*entities.InitResult, *dto.ErrorResponse)
	GenerateDeviceToken(body *svDto.SvDeviceInitRequest, serverUrl string) (string, error)
}

func NewCommonUsecase(repo repositories.CommonRepository) CommonUsecase {
	return &commonUsecase{repo: repo}
}

func (u *commonUsecase) GetConfig(req clDto.ConfigRequest) (*entities.Config, error) {
	// 대부분의 시스템에서는 단일 파일 다운로드 시 다음과 같은 패턴을 따릅니다:
	// Content-Type: application/octet-stream 또는 해당 파일 타입 (예: application/json, text/plain, application/x-yaml 등)

	// 설정으로 주입해야 하나?
	fileName := req.FileName
	filePath := filepath.Join("./config", fileName)

	// 파일 읽기
	content, err := os.ReadFile(filePath)
	fmt.Println("읽은 파일 ", content)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return &entities.Config{
		FileName: fileName,
		Content:  content,
	}, nil
}

func (u *commonUsecase) DeviceInit(body *svDto.SvDeviceInitRequest, ctx context.Context) (*entities.InitResult, *dto.ErrorResponse) {

	// DB 조회
	result, err := u.repo.GetConnectInfo(body.WorksCode)
	if err != nil {
		return &entities.InitResult{}, &dto.ErrorResponse{
			Code:    consts.E_102,
			Message: consts.E_102_MSG,
		}
	}

	// AUTH에 JWT 요청
	result.AppToken, err = u.GenerateDeviceToken(body, result.ConnectInfo)
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

func toWorksConfigEntity(worksCode string, device string) entities.GetWorksConfig {

	return entities.GetWorksConfig{
		WorksCode: worksCode,
		Device:    device,
	}

}

func (u *commonUsecase) GenerateDeviceToken(body *svDto.SvDeviceInitRequest, serverUrl string) (string, error) {
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
