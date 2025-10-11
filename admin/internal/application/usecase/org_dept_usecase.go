package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/domain/orgDept/entity"
	"admin/internal/domain/orgDept/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgDeptsUsecase struct {
	repository repository.OrgDeptRepository
	// serverUrl 추가하기
}

type OrgDeptsUsecase interface {
	RegisterDept(ctx context.Context, input input.RegisterDeptInput) (int, error)
	DeleteDept(ctx context.Context, input input.DeleteDeptInput) (int, error)
}

func NewOrgDeptsUsecase(repository repository.OrgDeptRepository) OrgDeptsUsecase {

	return &orgDeptsUsecase{
		repository: repository,
	}

}

func (r *orgDeptsUsecase) RegisterDept(ctx context.Context, input input.RegisterDeptInput) (int, error) {
	entity := entity.MakeRegisterDeptEntity(input.DeptCode, input.DeptOrg, input.ParentDeptCode, input.KoLang, input.EnLang, input.JpLang, input.ZhLang, input.RuLang, input.ZhLang, input.Header)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

var serverUrl = "172.16.10.114"

func createDepartmentInOrg(ctx context.Context, entity entity.RegisterDeptEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://" + serverUrl + "/org/server/v1/department"

	log.Println("org department create api 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 to 서버 인증 처리 필요

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

func (r *orgDeptsUsecase) DeleteDept(ctx context.Context, input input.DeleteDeptInput) (int, error) {
	entity := entity.MakeDeleteDeptEntity(input.DeptOrg, input.DeptCode)

	// 이 함수 내부에서 호출
	result, err := deleteDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func deleteDepartmentInOrg(ctx context.Context, entity entity.DeleteDeptEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://" + serverUrl + "/org/server/v1/department"

	log.Println("org department delete api 호출! url : ", url)

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
