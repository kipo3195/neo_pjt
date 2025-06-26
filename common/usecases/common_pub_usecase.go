package usecases

import (
	"bytes"
	clDto "common/dto/client"
	svDto "common/dto/server"
	"common/repositories"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type commonPubUsecase struct {
	repo repositories.CommonPubRepository
}

type CommonPubUsecase interface {
	AppValidation(ctx context.Context, body clDto.AppValidationRequest) (bool, error)
}

func NewCommonPubUsecase(repo repositories.CommonPubRepository) CommonPubUsecase {
	return &commonPubUsecase{
		repo: repo,
	}
}

func (r *commonPubUsecase) AppValidation(ctx context.Context, body clDto.AppValidationRequest) (bool, error) {

	//data, err := getAppTokenValidationInAuth(toAppTokenValidationRequest(body.AppToken, body.Uuid))

	// if err != nil {
	// 	// 에러 정의 후 response
	// }
	// if !data {
	// 	// 인증 실패, 에러 정의 후 response.
	// }

	// return flag, nil
	return true, nil
}

func toAppTokenValidationRequest(appToken string, uuid string) svDto.AppTokenValidationRequest {
	return svDto.AppTokenValidationRequest{
		AppToken: appToken,
		Uuid:     uuid,
	}

}

func getAppTokenValidationInAuth(dto svDto.AppTokenValidationRequest) (*svDto.AppTokenValidationResponse, error) {

	// JSON 변환
	jsonData, err := json.Marshal(dto)
	if err != nil {
		return &svDto.AppTokenValidationResponse{}, err
	}

	// POST 요청 보내기
	url := "http://" + "" + "/common/sv1/device-init"

	log.Println("auth service 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &svDto.AppTokenValidationResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // works 서버 호출시 필요한 키 작성하기 TODO

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &svDto.AppTokenValidationResponse{}, err
	}

	defer resp.Body.Close()

	var responseBody *svDto.AppTokenValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		fmt.Println("serverReponse 파싱시 에러")
		return &svDto.AppTokenValidationResponse{}, err
	}

	return responseBody, nil
}
