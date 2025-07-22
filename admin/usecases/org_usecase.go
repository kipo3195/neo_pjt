package usecases

import (
	clOrgReqDto "admin/dto/client/org/request"
	"admin/entities"
	orgEntity "admin/entities/org"
	"admin/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
	CreateDepartment(ctx context.Context, requestDTO clOrgReqDto.CreateDeptRequestDTO) (int, error)
	DeleteDepartment(ctx context.Context, requestDTO clOrgReqDto.DeleteDeptRequestDTO) (int, error)
	CreateOrgFile(ctx context.Context, requestDTO clOrgReqDto.CreateOrgFileRequestDTO) (int, error)

	CreateDeptUser(ctx context.Context, requestDTO clOrgReqDto.CreateDeptUserRequestDTO) (int, error)
	DeleteDeptUser(ctx context.Context, requestDTO clOrgReqDto.DeleteDeptUserRequestDTO) (int, error)
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) CreateDepartment(ctx context.Context, requestDTO clOrgReqDto.CreateDeptRequestDTO) (int, error) {
	entity := toCreateDepartmentEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateDepartmentEntity(body clOrgReqDto.CreateDeptRequestBody) orgEntity.CreateDepartmentEntity {

	return orgEntity.CreateDepartmentEntity{
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

func createDepartmentInOrg(ctx context.Context, entity orgEntity.CreateDepartmentEntity) (int, error) {

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

func (r *orgUsecase) DeleteDepartment(ctx context.Context, req clOrgReqDto.DeleteDeptRequestDTO) (int, error) {
	entity := toDeleteDepartmentEntity(req.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := deleteDepartmentInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toDeleteDepartmentEntity(body clOrgReqDto.DeleteDeptRequestBody) entities.DeleteDepartmentEntity {

	return entities.DeleteDepartmentEntity{
		DeptCode: body.DeptCode,
		DeptOrg:  body.DeptOrg,
	}
}

func deleteDepartmentInOrg(ctx context.Context, entity entities.DeleteDepartmentEntity) (int, error) {

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

func (r *orgUsecase) CreateOrgFile(ctx context.Context, requestDTO clOrgReqDto.CreateOrgFileRequestDTO) (int, error) {

	entity := toCreateOrgFileEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createOrgFileInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateOrgFileEntity(body clOrgReqDto.CreateOrgFileRequestBody) entities.CreateOrgFileEntity {
	return entities.CreateOrgFileEntity{
		OrgCode: body.OrgCode,
	}
}

func createOrgFileInOrg(ctx context.Context, entity entities.CreateOrgFileEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/org/file"

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

func (r *orgUsecase) CreateDeptUser(ctx context.Context, requestDTO clOrgReqDto.CreateDeptUserRequestDTO) (int, error) {

	entity := toCreateDeptUserEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateDeptUserEntity(body clOrgReqDto.CreateDeptUserRequestBody) entities.CreateDeptUserEntity {
	return entities.CreateDeptUserEntity{
		UserHash:             body.UserHash,
		DeptCode:             body.DeptCode,
		PositionCode:         body.PositionCode,
		RoleCode:             body.RoleCode,
		IsConcurrentPosition: body.IsConcurrentPosition,
	}
}

func createDeptUserInOrg(ctx context.Context, entity entities.CreateDeptUserEntity) (int, error) {

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

func (r *orgUsecase) DeleteDeptUser(ctx context.Context, requestDTO clOrgReqDto.DeleteDeptUserRequestDTO) (int, error) {
	entity := toDeleteDeptUserEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := deleteDeptUserInOrg(ctx, entity)

	if err != nil {
		return result, err
	}

	return result, nil
}

func toDeleteDeptUserEntity(body clOrgReqDto.DeleteDeptUserRequestBody) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: body.UserHash,
		DeptCode: body.DeptCode,
		DeptOrg:  body.DeptOrg,
	}

}

func deleteDeptUserInOrg(ctx context.Context, entity entities.DeleteDeptUserEntity) (int, error) {

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
