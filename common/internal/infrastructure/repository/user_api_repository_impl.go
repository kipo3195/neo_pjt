package repository

import (
	"bytes"
	"common/internal/domain/user/entity"
	"common/internal/domain/user/repository"
	"common/internal/infrastructure/dto/user"
	"common/pkg/dto"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type userAPIRepository struct {
	httpClient *http.Client
}

func NewUserAPIRepository() repository.UserAPIRepository {
	return &userAPIRepository{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (r *userAPIRepository) UserAuthRegistInAuth(ctx context.Context, id string, entity entity.UserRegisterInfoEntity, challenge string) (string, error) {

	url := "http://" + "" + "/auth/server/v1/user/auth/info/register"
	log.Println("auth service 호출! url : ", url)

	userAuth := make([]user.UserAuthRegisterDto, 0)

	userAuth = append(userAuth, user.UserAuthRegisterDto{
		UserId:   id,
		Salt:     entity.Salt,
		UserHash: challenge,
		UserAuth: entity.Hash,
	})

	reqBody := user.UserAuthRegisterRequest{
		UserAuth: userAuth,
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
	var result dto.ServerResponseDTO[*user.UserAuthRegisterResponse]

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("serverReponse 파싱시 에러")
		return "", err
	}

	return result.Data.Result, nil
}
