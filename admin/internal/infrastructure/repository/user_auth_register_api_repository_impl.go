package repository

import (
	"admin/internal/domain/serviceUser/entity"
	"admin/internal/domain/serviceUser/repository"
	"admin/pkg/dto"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type userAuthRegisterApiRepositoryImpl struct {
}

func (r *userAuthRegisterApiRepositoryImpl) NewServiceUserApiRepository() repository.UserAuthRegisterApiRepository {

	return &userAuthRegisterApiRepositoryImpl{}

}

func (r *userAuthRegisterApiRepositoryImpl) UserAuthRegisterInAuth(ctx context.Context, entity []entity.ServiceUserEntity) error {

	url := "http://" + "" + "/auth/server/v1/user/auth/info/register"
	log.Println("auth service 호출! url : ", url)

	userAuth := make([]user.UserAuthRegisterDto, 0)

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
