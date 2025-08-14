package client

import (
	"bytes"
	"common/internal/domains/appToken/dto/client/requestDTO"
	authResponseDTO "common/internal/domains/appToken/dto/external/authResponseDTO"
	repositories "common/internal/domains/appToken/repositories/client"
	commonConsts "common/pkg/consts"
	"common/pkg/dto"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type appTokenUsecase struct {
	repository repositories.AppTokenRepository
}

type AppTokenUsecase interface {
	AppTokenReIssue(ctx context.Context, requestDTO requestDTO.AppTokenRefreshRequestDTO) (*authResponseDTO.AppTokenRefreshResponseBody, error)
}

func NewAppTokenUsecase(repository repositories.AppTokenRepository) AppTokenUsecase {
	return &appTokenUsecase{
		repository: repository,
	}
}

func (r *appTokenUsecase) AppTokenReIssue(ctx context.Context, requestDTO requestDTO.AppTokenRefreshRequestDTO) (*authResponseDTO.AppTokenRefreshResponseBody, error) {

	// marshal
	requestBody, err := json.Marshal(requestDTO.Body)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", commonConsts.ErrServerError)
	}

	url := "http://auth-service/auth/sv1/app-token-refresh"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", commonConsts.ErrServerError)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth call failed: %w", commonConsts.ErrServerError)
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

	var result dto.ServerResponseDTO[*authResponseDTO.AppTokenRefreshResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	responseDTO := result.Data

	log.Println("auth service 호출 end !")
	return &responseDTO.Body, nil

}
