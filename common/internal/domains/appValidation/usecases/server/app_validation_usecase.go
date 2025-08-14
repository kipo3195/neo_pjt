package server

import (
	"bytes"
	"common/internal/consts"
	"common/internal/domains/appValidation/dto/external/authRequestDTO"
	"common/internal/domains/appValidation/dto/server/requestDTO"
	repositories "common/internal/domains/appValidation/repositories/server"
	"common/internal/infra/storage"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type appValidationUsecase struct {
	repository        repositories.AppValidationRepository
	configHashStorage storage.ConfigHashStorage
}

type AppValidationUsecase interface {
	AppValidation(ctx context.Context, requestDTO requestDTO.AppValidationRequestDTO) (bool, error)
}

func NewAppValidationUsecase(repository repositories.AppValidationRepository, configHashStorage storage.ConfigHashStorage) AppValidationUsecase {
	return &appValidationUsecase{
		repository:        repository,
		configHashStorage: configHashStorage,
	}
}

func (r *appValidationUsecase) AppValidation(ctx context.Context, requestDTO requestDTO.AppValidationRequestDTO) (bool, error) {

	//http statuscode를 리턴함.

	var tokenType, token string

	if requestDTO.Body.AppToken != "" {
		tokenType = "appToken"
		token = requestDTO.Body.AppToken
	} else if requestDTO.Body.AccessToken != "" {
		tokenType = "accessToken"
		token = requestDTO.Body.AccessToken
	} else {
		log.Println("token type error.")
		return false, consts.ErrServerError
	}
	_, err := getAppTokenValidationInAuth(ctx, toAppTokenValidationRequest(token, tokenType, requestDTO.Body.Uuid))

	if err != nil {
		// 에러 정의 후 response
		return false, consts.ErrServerError
	}

	return true, nil
}

func toAppTokenValidationRequest(token string, tokenType string, uuid string) authRequestDTO.AppTokenValidationRequestDTO {
	return authRequestDTO.AppTokenValidationRequestDTO{
		Body: authRequestDTO.AppTokenValidationRequestBody{
			Token:     token,
			TokenType: tokenType,
			Uuid:      uuid,
		},
	}
}

func getAppTokenValidationInAuth(ctx context.Context, requestDTO authRequestDTO.AppTokenValidationRequestDTO) (int, error) {

	// POST 요청 보내기
	url := "http://172.16.10.114/auth/sv1/app-token-validation"
	log.Println("auth service 호출! url : ", url)

	// JSON 변환
	serverRequestBody, err := json.Marshal(requestDTO.Body)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(serverRequestBody))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err

	}

	defer resp.Body.Close()

	if err != nil {
		log.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return http.StatusInternalServerError, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return http.StatusInternalServerError, fmt.Errorf("request failed: %w", err)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}
