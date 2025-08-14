package client

import (
	"admin/internal/domains/orgFile/dto/client/requestDTO"
	"admin/internal/domains/orgFile/entities"
	repository "admin/internal/domains/orgFile/repositories/client"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type orgFileUsecase struct {
	repository repository.OrgFileRepository
}

type OrgFileUsecase interface {
	CreateOrgFile(ctx context.Context, requestDTO requestDTO.CreateOrgFileRequestDTO) (int, error)
}

func NewOrgFileUsecase(repository repository.OrgFileRepository) OrgFileUsecase {

	return &orgFileUsecase{
		repository: repository,
	}
}

func (r *orgFileUsecase) CreateOrgFile(ctx context.Context, requestDTO requestDTO.CreateOrgFileRequestDTO) (int, error) {

	entity := toCreateOrgFileEntity(requestDTO.Body)
	// org 서비스 호출

	// 이 함수 내부에서 호출
	result, err := createOrgFileInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func toCreateOrgFileEntity(body requestDTO.CreateOrgFileRequestBody) entities.CreateOrgFileEntity {
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
