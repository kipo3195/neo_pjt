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
	repository    repositories.AppValidationRepository
	configStorage storage.ConfigStorage
}

type AppValidationUsecase interface {
	AppValidation(ctx context.Context, requestDTO requestDTO.AppValidationRequestDTO) (bool, error)
}

func NewAppValidationUsecase(repository repositories.AppValidationRepository, configStorage storage.ConfigStorage) AppValidationUsecase {
	return &appValidationUsecase{
		repository:    repository,
		configStorage: configStorage,
	}
}

func (r *appValidationUsecase) AppValidation(ctx context.Context, requestDTO requestDTO.AppValidationRequestDTO) (bool, error) {

	//http statuscode를 리턴함.
	_, err := getAppTokenValidationInAuth(ctx, toAppTokenValidationRequest(requestDTO.Body))

	if err != nil {
		// 에러 정의 후 response
		return false, consts.ErrServerError
	}

	// skin 검증
	skinHash, err := r.configStorage.GetHash(consts.SKIN)
	if err != nil {
		log.Println("서버에 skin hash정보가 없음.")
		return false, consts.ErrSkinHashInvalid
	}

	if skinHash != requestDTO.Body.SkinHash {
		log.Println("서버의 skin hash 정보와 다름 server skin hash : ", skinHash)
		return false, consts.ErrSkinHashInvalid
	}

	// config 검증
	configHash, err := r.configStorage.GetHash(consts.CONFIG)
	if err != nil {
		log.Println("서버에 config hash정보가 없음.")
		return false, consts.ErrConfigHashInvalid
	}

	if configHash != requestDTO.Body.ConfigHash {
		log.Println("서버의 config hash 정보와 다름 server config hash : ", configHash)
		return false, consts.ErrConfigHashInvalid
	}

	return true, nil
}

func toAppTokenValidationRequest(body requestDTO.AppValidationRequestBody) authRequestDTO.AppTokenValidationRequestDTO {
	return authRequestDTO.AppTokenValidationRequestDTO{
		Body: authRequestDTO.AppTokenValidationRequestBody{
			AppToken: body.AppToken,
			Uuid:     body.Uuid,
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
