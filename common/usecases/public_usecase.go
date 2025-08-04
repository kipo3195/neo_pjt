package usecases

import (
	"bytes"
	"common/consts"
	clCommonReqDto "common/dto/client/request"
	dto "common/dto/common"
	svAuthResDto "common/dto/server/auth/response"
	"common/infra/storage"
	"common/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type publicUsecase struct {
	repo          repositories.PublicRepository
	configStorage storage.ConfigStorage
}

type PublicUsecase interface {
	AppTokenReIssue(ctx context.Context, requestDTO clCommonReqDto.AppTokenRefreshRequestDTO) (*svAuthResDto.AppTokenRefreshResponseBody, error)
}

func NewPublicUsecase(repo repositories.PublicRepository, configStorage storage.ConfigStorage) PublicUsecase {
	return &publicUsecase{
		repo:          repo,
		configStorage: configStorage,
	}
}

func (r *publicUsecase) AppTokenReIssue(ctx context.Context, requestDTO clCommonReqDto.AppTokenRefreshRequestDTO) (*svAuthResDto.AppTokenRefreshResponseBody, error) {

	// marshal
	requestBody, err := json.Marshal(requestDTO.Body)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", consts.ErrServerError)
	}

	url := "http://auth-service/auth/sv1/app-token-refresh"
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

	var result dto.ServerResponseDTO[*svAuthResDto.AppTokenRefreshResponseDTO]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return nil, err
	}

	responseDTO := result.Data

	log.Println("auth service 호출 end !")
	return &responseDTO.Body, nil

}
