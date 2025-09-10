package usecase

import (
	"bytes"
	"common/internal/consts"
	"common/internal/delivery/dto/appToken"
	"common/internal/domain/appToken/entity"
	"common/internal/domain/appToken/repository"
	"common/pkg/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type appTokenUsecase struct {
	repository repository.AppTokenRepository
}

type AppTokenUsecase interface {
	AppTokenReIssueInAuth(ctx context.Context, requestDTO appToken.AppTokenRefreshRequestDTO) (*appToken.AppTokenRefreshResponseBody, error)
	GenerateAppTokenInAuth(requestDTO *appToken.GenerateAppTokenRequestDTO, serverUrl string) (*entity.GenerateAppToken, error)
}

func NewAppTokenUsecase(repository repository.AppTokenRepository) AppTokenUsecase {
	return &appTokenUsecase{
		repository: repository,
	}
}

func (r *appTokenUsecase) AppTokenReIssueInAuth(ctx context.Context, requestDTO appToken.AppTokenRefreshRequestDTO) (*appToken.AppTokenRefreshResponseBody, error) {

	// marshal
	requestBody, err := json.Marshal(requestDTO.Body)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", consts.ErrServerError)
	}

	url := "http://auth-service/auth/server/v1/app-token-refresh"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", consts.ErrServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth call failed: %w", consts.ErrServerError)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("request failed: %w", err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var result dto.ServerResponseDTO[*appToken.AppTokenRefreshResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	responseDTO := result.Data

	log.Println("auth service 호출 end !")
	return &responseDTO.Body, nil

}

func (r *appTokenUsecase) GenerateAppTokenInAuth(requestDTO *appToken.GenerateAppTokenRequestDTO, serverUrl string) (*entity.GenerateAppToken, error) {
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
	var result dto.ServerResponseDTO[*appToken.GenerateAppTokenResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	return toGenerateAppTokenEntity(result.Data), nil
}

func toGenerateAppTokenEntity(data *appToken.GenerateAppTokenResponseDTO) *entity.GenerateAppToken {

	return &entity.GenerateAppToken{
		AppToken:     data.Body.AppToken,
		RefreshToken: data.Body.RefreshToken,
	}
}
