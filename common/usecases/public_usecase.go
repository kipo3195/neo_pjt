package usecases

import (
	"bytes"
	"common/consts"
	clDto "common/dto/client"
	authDto "common/dto/server/auth"
	"common/entities"
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
	AppValidation(ctx context.Context, body clDto.AppValidationRequest) (bool, error)
}

func NewPublicUsecase(repo repositories.PublicRepository, configStorage storage.ConfigStorage) PublicUsecase {
	return &publicUsecase{
		repo:          repo,
		configStorage: configStorage,
	}
}

func (r *publicUsecase) AppValidation(ctx context.Context, body clDto.AppValidationRequest) (bool, error) {

	data, err := getAppTokenValidationInAuth(toAppTokenValidationRequest(body.AppToken, body.Uuid))

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
		fmt.Println("서버에 skin hash정보가 없음.")
		return false, consts.ErrSkinHashInvalid
	}

	if skinHash != body.SkinHash {
		fmt.Println("서버의 skin hash 정보와 다름 server skin hash : ", skinHash)
		return false, consts.ErrSkinHashInvalid
	}

	// config 검증
	configHash, err := r.configStorage.GetHash(consts.CONFIG)
	if err != nil {
		fmt.Println("서버에 config hash정보가 없음.")
		return false, consts.ErrConfigHashInvalid
	}

	if configHash != body.ConfigHash {
		fmt.Println("서버의 config hash 정보와 다름 server config hash : ", configHash)
		return false, consts.ErrConfigHashInvalid
	}

	return true, nil
}

func toAppTokenValidationRequest(appToken string, uuid string) authDto.AppTokenValidationRequest {
	return authDto.AppTokenValidationRequest{
		AppToken: appToken,
		Uuid:     uuid,
	}

}

func getAppTokenValidationInAuth(dto authDto.AppTokenValidationRequest) (entities.AppTokenValitaionResponseEntity, error) {

	// JSON 변환
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return entities.AppTokenValitaionResponseEntity{}, err
	}

	// POST 요청 보내기
	url := "http://172.16.10.114/auth/sv1/app-token-validation"

	log.Println("auth service 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
		fmt.Println("serverReponse 파싱시 에러")
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
