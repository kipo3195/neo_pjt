package client

import (
	"admin/internal/domains/orgDeptUsers/dto/client/requestDTO"
	entities "admin/internal/domains/orgDeptUsers/entities"
	clientRepository "admin/internal/domains/orgDeptUsers/repositories/client"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgDeptUsersUsecase struct {
	repository clientRepository.OrgDeptUsersRepository
}

type OrgDeptUsersUsecase interface {
	CreateDeptUser(ctx context.Context, requestDTO requestDTO.CreateDeptUserRequestDTO) (int, error)
	DeleteDeptUser(ctx context.Context, requestDTO requestDTO.DeleteDeptUserRequestDTO) (int, error)
}

func NewOrgDeptUsersUsecase(repository clientRepository.OrgDeptUsersRepository) OrgDeptUsersUsecase {

	return &orgDeptUsersUsecase{
		repository: repository,
	}

}

func (r *orgDeptUsersUsecase) DeleteDeptUser(ctx context.Context, requestDTO requestDTO.DeleteDeptUserRequestDTO) (int, error) {
	entity := toDeleteDeptUserEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := deleteDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func toDeleteDeptUserEntity(body requestDTO.DeleteDeptUserRequestBody) entities.DeleteDeptUsersEntity {

	return entities.DeleteDeptUsersEntity{
		UserHash: body.UserHash,
		DeptCode: body.DeptCode,
		DeptOrg:  body.DeptOrg,
	}

}

func deleteDeptUserInOrg(ctx context.Context, entity entities.DeleteDeptUsersEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments/user"

	log.Println("create org file api 호출! url : ", url)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return http.StatusInternalServerError, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return http.StatusInternalServerError, fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}

func (r *orgDeptUsersUsecase) CreateDeptUser(ctx context.Context, requestDTO requestDTO.CreateDeptUserRequestDTO) (int, error) {

	entity := toCreateDeptUserEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateDeptUserEntity(body requestDTO.CreateDeptUserRequestBody) entities.RegisterDeptUsersEntity {
	return entities.RegisterDeptUsersEntity{
		UserHash:             body.UserHash,
		DeptCode:             body.DeptCode,
		PositionCode:         body.PositionCode,
		RoleCode:             body.RoleCode,
		IsConcurrentPosition: body.IsConcurrentPosition,
	}
}

func createDeptUserInOrg(ctx context.Context, entity entities.RegisterDeptUsersEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments/user"

	log.Println("create org file api 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return http.StatusInternalServerError, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return http.StatusInternalServerError, fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.StatusCode, fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}
