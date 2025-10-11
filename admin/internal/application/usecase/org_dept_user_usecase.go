package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/domain/orgDeptUser/entity"
	"admin/internal/domain/orgDeptUser/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgDeptUsersUsecase struct {
	repository repository.OrgDeptUserRepository
}

type OrgDeptUsersUsecase interface {
	RegistDeptUser(ctx context.Context, input input.RegistDeptUserInput) (int, error)
	DeleteDeptUser(ctx context.Context, input input.DeleteDeptUserInput) (int, error)
}

func NewOrgDeptUsersUsecase(repository repository.OrgDeptUserRepository) OrgDeptUsersUsecase {

	return &orgDeptUsersUsecase{
		repository: repository,
	}

}

func (r *orgDeptUsersUsecase) DeleteDeptUser(ctx context.Context, input input.DeleteDeptUserInput) (int, error) {
	entity := entity.MakeDeleteDeptUserEntity(input.UserHash, input.DeptOrg, input.DeptCode)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := deleteDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func deleteDeptUserInOrg(ctx context.Context, entity entity.DeleteDeptUserEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://" + serverUrl + "/org/server/v1/department/user"

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

func (r *orgDeptUsersUsecase) RegistDeptUser(ctx context.Context, input input.RegistDeptUserInput) (int, error) {

	entity := entity.MakeRegistDeptUserEntity(input.UserHash, input.DeptOrg, input.DeptCode, input.RoleCode, input.PositionCode, input.IsConcurrentPosition)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := registDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func registDeptUserInOrg(ctx context.Context, entity entity.RegistDeptUserEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://" + serverUrl + "/org/sv1/department/user"

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
