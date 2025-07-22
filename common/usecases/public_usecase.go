package usecases

import (
	"bytes"
	"common/consts"
	clDto "common/dto/client"
	clCommonReqDto "common/dto/client/request"
	dto "common/dto/common"
	authDto "common/dto/server/auth"
	"common/entities"
	"common/infra/storage"
	"common/repositories"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type publicUsecase struct {
	repo          repositories.PublicRepository
	configStorage storage.ConfigStorage
}

type PublicUsecase interface {
	AppValidation(ctx context.Context, requestDTO clCommonReqDto.AppValidationRequestDTO) (bool, error)
	AppTokenReIssue(ctx context.Context, body clDto.AppTokenRefreshRequest) (*authDto.AppTokenRefreshResponse, error)
}

func NewPublicUsecase(repo repositories.PublicRepository, configStorage storage.ConfigStorage) PublicUsecase {
	return &publicUsecase{
		repo:          repo,
		configStorage: configStorage,
	}
}

func (r *publicUsecase) AppValidation(ctx context.Context, requestDTO clCommonReqDto.AppValidationRequestDTO) (bool, error) {

	data, err := getAppTokenValidationInAuth(toAppTokenValidationRequest(requestDTO.Body))

	if err != nil {
		// 에러 정의 후 response
		return false, consts.ErrServerError
	}

	if data.Result != "success" {
		// 인증 실패, 에러 정의 후 response.
		return false, consts.ErrInvalidClaims
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

func toAppTokenValidationRequest(body clCommonReqDto.AppValidationRequestBody) authDto.AppTokenValidationRequest {
	return authDto.AppTokenValidationRequest{
		AppToken: body.AppToken,
		Uuid:     body.Uuid,
	}

}

func getAppTokenValidationInAuth(body authDto.AppTokenValidationRequest) (entities.AppTokenValitaionResponseEntity, error) {

	// JSON 변환
	serverRequestBody, err := json.Marshal(body)
	if err != nil {
		return entities.AppTokenValitaionResponseEntity{}, err
	}

	// POST 요청 보내기
	url := "http://172.16.10.114/auth/sv1/app-token-validation"

	log.Println("auth service 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(serverRequestBody))
	if err != nil {
		return entities.AppTokenValitaionResponseEntity{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return entities.AppTokenValitaionResponseEntity{}, err
	}

	defer resp.Body.Close()

	var responseBody authDto.AppTokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return entities.AppTokenValitaionResponseEntity{}, err
	}

	return toAppTokenValidationResponseEntity(responseBody), nil
}

func toAppTokenValidationResponseEntity(dto authDto.AppTokenValidationResponse) entities.AppTokenValitaionResponseEntity {
	return entities.AppTokenValitaionResponseEntity{
		Result: dto.Result,
		Data:   dto.Data,
	}
}

func (r *publicUsecase) AppTokenReIssue(ctx context.Context, body clDto.AppTokenRefreshRequest) (*authDto.AppTokenRefreshResponse, error) {

	// marshal
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %w", consts.ErrServerError)
	}

	url := "http://auth-service/auth/sv1/app-token-refresh"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", consts.ErrServerError)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("auth call failed: %w", consts.ErrServerError)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body failed: %w", consts.ErrServerError)
	}

	// 응답 코드에 따라 분기
	switch resp.StatusCode {
	case http.StatusOK:
		var result dto.ServerResponse[*authDto.AppTokenRefreshResponse]
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", consts.ErrServerError)
		}
		return result.Data, nil

	case http.StatusBadRequest:
		var errResp dto.ErrorResponse
		if err := json.Unmarshal(bodyBytes, &errResp); err != nil {
			return nil, fmt.Errorf("unmarshal error body failed: %w", consts.ErrServerError)
		}

		switch errResp.Code {
		case "AUTH_F001":
			return nil, consts.ErrRefreshTokenAuthInvalid
		case "AUTH_F002":
			return nil, consts.ErrRefreshTokenAuthExpired
		default:
			return nil, fmt.Errorf("unknown auth error: %w", consts.ErrServerError)
		}

	default:
		return nil, consts.ErrServerError
	}
}
