package usecases

import (
	clDto "admin/dto/client"
	"admin/entities"
	"admin/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type adminOrgUsecase struct {
	repo repositories.AdminOrgRepository
}

type AdminOrgUsecase interface {
	CreateDepartment(ctx context.Context, req clDto.CreateDeptRequest) (interface{}, error)
	DeleteDepartment(ctx context.Context, req clDto.DeleteDeptRequest) (interface{}, error)
	CreateOrgFile(ctx context.Context, req clDto.CreateOrgFileRequest) (interface{}, error)

	CreateDeptUser(ctx context.Context, req clDto.CreateDeptUserRequest) (interface{}, error)
	DeleteDeptUser(ctx context.Context, req clDto.DeleteDeptUserRequest) (interface{}, error)
}

func NewAdminOrgUsecase(repo repositories.AdminOrgRepository) AdminOrgUsecase {
	return &adminOrgUsecase{repo: repo}
}

func (r *adminOrgUsecase) CreateDepartment(ctx context.Context, req clDto.CreateDeptRequest) (interface{}, error) {
	entity := toCreateDepartmentEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	err := createDepartmentInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func toCreateDepartmentEntity(req clDto.CreateDeptRequest) entities.CreateDepartmentEntity {

	return entities.CreateDepartmentEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		KoLang:         req.KoLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		ZhLang:         req.ZhLang,
		ViLang:         req.ViLang,
		Header:         req.Header,
	}
}

func createDepartmentInOrg(ctx context.Context, entity entities.CreateDepartmentEntity) error {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments"

	fmt.Println("org department create api 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken") // 서버 to 서버 인증 처리 필요

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return nil
}

func (r *adminOrgUsecase) DeleteDepartment(ctx context.Context, req clDto.DeleteDeptRequest) (interface{}, error) {
	entity := toDeleteDepartmentEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	err := deleteDepartmentInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	// 여기 수정해야할듯.
	return nil, nil
}

func toDeleteDepartmentEntity(req clDto.DeleteDeptRequest) entities.DeleteDepartmentEntity {

	return entities.DeleteDepartmentEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func deleteDepartmentInOrg(ctx context.Context, entity entities.DeleteDepartmentEntity) error {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments"

	fmt.Println("org department delete api 호출! url : ", url)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return nil
}

func (r *adminOrgUsecase) CreateOrgFile(ctx context.Context, req clDto.CreateOrgFileRequest) (interface{}, error) {

	entity := toCreateOrgFileEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := CreateOrgFileInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func toCreateOrgFileEntity(req clDto.CreateOrgFileRequest) entities.CreateOrgFileEntity {
	return entities.CreateOrgFileEntity{
		OrgCode: req.OrgCode,
	}
}

func CreateOrgFileInOrg(ctx context.Context, entity entities.CreateOrgFileEntity) (interface{}, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/org/file"

	fmt.Println("create org file api 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}

func (r *adminOrgUsecase) CreateDeptUser(ctx context.Context, req clDto.CreateDeptUserRequest) (interface{}, error) {

	entity := toCreateDeptUserEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := CreateDeptUserInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func toCreateDeptUserEntity(req clDto.CreateDeptUserRequest) entities.CreateDeptUserEntity {
	return entities.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
	}
}

func CreateDeptUserInOrg(ctx context.Context, entity entities.CreateDeptUserEntity) (interface{}, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments/user"

	fmt.Println("create org file api 호출! url : ", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}

func (r *adminOrgUsecase) DeleteDeptUser(ctx context.Context, req clDto.DeleteDeptUserRequest) (interface{}, error) {
	entity := toDeleteDeptUserEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := DeptUserInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func toDeleteDeptUserEntity(req clDto.DeleteDeptUserRequest) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}

}

func DeptUserInOrg(ctx context.Context, entity entities.DeleteDeptUserEntity) (interface{}, error) {

	payload, _ := json.Marshal(entity)

	url := "http://172.16.10.114/org/sv1/departments/user"

	fmt.Println("create org file api 호출! url : ", url)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer serverToken")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("org error : ", err)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("request cancelled or timed out: %w", ctx.Err())
		default:
			return nil, fmt.Errorf("request failed: %w", err)
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("org service returned status %d", resp.StatusCode)
	}

	return http.StatusOK, nil
}
