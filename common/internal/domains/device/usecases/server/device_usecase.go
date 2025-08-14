package server

import (
	"bytes"
	"common/internal/domains/device/dto/external/authResponseDTO"
	"common/internal/domains/device/dto/server/requestDTO"
	"common/internal/domains/device/entities"
	"common/internal/domains/device/repositories/serverRepository"
	"common/pkg/consts"
	"common/pkg/dto"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type deviceUsecase struct {
	repository serverRepository.DeviceRepository
}

type DeviceUsecase interface {
	DeviceInit(ctx context.Context, body *requestDTO.DeviceInitRequest) (*entities.InitResult, error)
}

func NewDeviceUsecase(repository serverRepository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{
		repository: repository,
	}
}

func (u *deviceUsecase) DeviceInit(ctx context.Context, body *requestDTO.DeviceInitRequest) (*entities.InitResult, error) {

	// DB 조회 connectInfo(접속 url)은 관리해야할 필요있음. 최초 이후에 클라이언트가 정보가 필요할때를 대비해서.
	connectInfo, err := u.repository.GetConnectInfo(body.WorksCode)
	if err != nil {
		return nil, consts.ErrDB
	}

	// AUTH에 JWT 요청
	issuedAppToken, err := generateAppToken(body, connectInfo.ServerUrl) //  serverUrl은 이후에 .env 또는 k8s의 secrets에서 읽기
	if err != nil {
		return nil, consts.ErrServerError
	}

	// 타임존, 언어, 앱 별 스킨 정보, 설정 정보 - GetConnectInfo와 합쳐서 트랜잭션 처리
	worksConfig, err := u.repository.GetWorksConfig(toWorksConfigEntity(body.WorksCode, body.Device), ctx)
	if err != nil {
		return nil, consts.ErrDB
	}

	log.Println("connectInfo:", connectInfo)
	log.Println("issuedAppToken:", issuedAppToken)
	log.Println("worksConfig ", worksConfig)

	return &entities.InitResult{
		ConnectInfo:    connectInfo,
		IssuedAppToken: issuedAppToken,
		WorksConfig:    worksConfig,
	}, nil
}

func toWorksConfigEntity(worksCode string, device string) entities.GetWorksConfig {

	return entities.GetWorksConfig{
		WorksCode: worksCode,
		Device:    device,
	}

}

func generateAppToken(body *requestDTO.DeviceInitRequest, serverUrl string) (*entities.IssuedAppToken, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": body.Uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Println("auth service 호출! 1")

	url := "http://" + serverUrl + "/auth/sv1/generate-app-token"
	//url := domain + "/auth/v1/generate-device-token"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 api key

	// POST 요청 보내기
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("auth service 호출 에러 1")
		return nil, err
	}
	defer resp.Body.Close()

	// 응답 출력
	var result dto.ServerResponseDTO[*authResponseDTO.DeviceInitResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	return toIssuedAppTokenEntity(result.Data), nil
}

func toIssuedAppTokenEntity(dto *authResponseDTO.DeviceInitResponse) *entities.IssuedAppToken {

	return &entities.IssuedAppToken{
		AppToken:     dto.AppToken,
		RefreshToken: dto.RefreshToken,
	}
}
