package client

import (
	"bytes"
	"common/internal/domains/appToken/dto/external/authResponseDTO"
	"common/internal/domains/appToken/dto/server/requestDTO"
	"common/internal/domains/appToken/entities"
	"common/pkg/dto"
	"encoding/json"
	"log"
	"net/http"
)

type appTokenUsecase struct {
}

type AppTokenUsecase interface {
	GenerateAppTokenInAuth(requestDTO *requestDTO.GenerateAppTokenRequestDTO, serverUrl string) (*entities.GenerateAppToken, error)
}

func NewAppTokenUsecase() AppTokenUsecase {
	return &appTokenUsecase{}
}

func (r *appTokenUsecase) GenerateAppTokenInAuth(requestDTO *requestDTO.GenerateAppTokenRequestDTO, serverUrl string) (*entities.GenerateAppToken, error) {
	// 소스 모듈화 처리하기
	data := map[string]string{
		"uuid": requestDTO.Body.Uuid,
	}

	// JSON 변환
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	log.Println("auth service 호출! 1")

	url := "http://auth-service/auth/server/v1/generate-app-token"
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
	var result dto.ServerResponseDTO[*authResponseDTO.GenerateAppTokenResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	return toGenerateAppTokenEntity(result.Data), nil
}

func toGenerateAppTokenEntity(data *authResponseDTO.GenerateAppTokenResponseDTO) *entities.GenerateAppToken {

	return &entities.GenerateAppToken{
		AppToken:     data.Body.AppToken,
		RefreshToken: data.Body.RefreshToken,
	}
}
