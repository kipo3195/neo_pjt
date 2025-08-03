package client

import (
	"admin/internal/domains/orgDepts/dto/client/requestDTO"
	entities "admin/internal/domains/orgDepts/entities"
	clientRepository "admin/internal/domains/orgDepts/repositories/client"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgDeptsUsecase struct {
	repository clientRepository.OrgDeptsRepository
}

type OrgDeptsUsecase interface {
	CreateDepartment(ctx context.Context, requestDTO requestDTO.RegisterDeptRequestDTO) (int, error)
	DeleteDepartment(ctx context.Context, requestDTO requestDTO.DeleteDeptRequestDTO) (int, error)
}

func NewOrgDeptsUsecase(repository clientRepository.OrgDeptsRepository) OrgDeptsUsecase {

	return &orgDeptsUsecase{
		repository: repository,
	}

}

func (r *orgDeptsUsecase) CreateDepartment(ctx context.Context, requestDTO requestDTO.RegisterDeptRequestDTO) (int, error) {
	entity := toCreateDepartmentEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateDepartmentEntity(body requestDTO.RegisterDeptRequestBody) entities.RegisterDeptsEntity {

	return entities.RegisterDeptsEntity{
		DeptCode:       body.DeptCode,
		DeptOrg:        body.DeptOrg,
		ParentDeptCode: body.ParentDeptCode,
		KoLang:         body.KoLang,
		EnLang:         body.EnLang,
		JpLang:         body.JpLang,
		ZhLang:         body.ZhLang,
		ViLang:         body.ViLang,
		Header:         body.Header,
	}
}

func createDepartmentInOrg(ctx context.Context, entity entities.RegisterDeptsEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments"

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

func (r *orgDeptsUsecase) DeleteDepartment(ctx context.Context, req requestDTO.DeleteDeptRequestDTO) (int, error) {
	entity := toDeleteDepartmentEntity(req.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := deleteDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toDeleteDepartmentEntity(body requestDTO.DeleteDeptRequestBody) entities.DeleteDeptsEntity {

	return entities.DeleteDeptsEntity{
		DeptCode: body.DeptCode,
		DeptOrg:  body.DeptOrg,
	}
}

func deleteDepartmentInOrg(ctx context.Context, entity entities.DeleteDeptsEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments"

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
