package usecases

import (
	"admin/dto"
	"admin/entities"
	"admin/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type adminUsecase struct {
	repo repositories.AdminRepository
}

type AdminUsecase interface {
	CreateDepartment(ctx context.Context, req dto.CreateDeptRequest) (interface{}, error)
	DeleteDepartment(ctx context.Context, req dto.DeleteDeptRequest) (interface{}, error)
}

func NewAdminUsecase(repo repositories.AdminRepository) AdminUsecase {
	return &adminUsecase{repo: repo}
}

func (r *adminUsecase) CreateDepartment(ctx context.Context, req dto.CreateDeptRequest) (interface{}, error) {
	entity := toCreateDepartmentEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	err := createDepartmentInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func toCreateDepartmentEntity(req dto.CreateDeptRequest) entities.CreateDepartmentEntity {

	return entities.CreateDepartmentEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		KrLang:         req.KrLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		CnLang:         req.CnLang,
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

func (r *adminUsecase) DeleteDepartment(ctx context.Context, req dto.DeleteDeptRequest) (interface{}, error) {
	entity := toDeleteDepartmentEntity(req)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	err := deleteDepartmentInOrg(ctx, entity)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func toDeleteDepartmentEntity(req dto.DeleteDeptRequest) entities.DeleteDepartmentEntity {

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
