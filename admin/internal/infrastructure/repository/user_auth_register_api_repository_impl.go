package repository

import (
	"admin/internal/domain/userAuthRegister/entity"
	"admin/internal/domain/userAuthRegister/repository"
	"admin/internal/infrastructure/dto/userAuth"
	"admin/pkg/dto"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type userAuthRegisterApiRepositoryImpl struct {
	domain     string
	httpClient *http.Client
}

func NewServiceUserApiRepository(domain string) repository.UserAuthRegisterApiRepository {

	return &userAuthRegisterApiRepositoryImpl{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		domain:     domain,
	}

}

func (r *userAuthRegisterApiRepositoryImpl) UserAuthRegisterInAuth(ctx context.Context, entity []entity.UserAuthRegisterEntity) (string, error) {

	url := "http://" + r.domain + "/auth/server/v1/user/auth/info/register"
	log.Println("auth service 호출! url : ", url)

	body := make([]userAuth.UserAuthRegisterDto, 0)

	for _, ua := range entity {

		temp := userAuth.UserAuthRegisterDto{
			UserId:   ua.UserId,
			UserHash: ua.UserHash,
			UserAuth: ua.UserAuth,
			Salt:     ua.Salt,
		}

		body = append(body, temp)
	}

	reqBody := userAuth.UserAuthRegisterRequest{
		UserAuth: body,
	}

	// 직렬화
	bodyByte, err := json.Marshal(reqBody)
	if err != nil {
		// 에러 처리
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "") // works 서버 호출시 필요한 키 작성하기 TODO

	client := r.httpClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// serverResponse로 전달받기 -> dto 뽑아내기 제네릭
	var result dto.ServerResponseDTO[*userAuth.UserAuthRegisterResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return "", err
	}

	return result.Data.Result, nil

}
