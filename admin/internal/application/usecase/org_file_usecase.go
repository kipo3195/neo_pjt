package usecase

import (
	"admin/internal/application/usecase/input"
	"admin/internal/domain/orgFile/entity"
	"admin/internal/domain/orgFile/repository"
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
	CreateOrgFile(ctx context.Context, input input.CreateOrgFileInput) (int, error)
}

func NewOrgFileUsecase(repository repository.OrgFileRepository) OrgFileUsecase {

	return &orgFileUsecase{
		repository: repository,
	}
}

func (r *orgFileUsecase) CreateOrgFile(ctx context.Context, input input.CreateOrgFileInput) (int, error) {

	entity := entity.MakeCreateOrgFileEntity(input.OrgCode)

	result, err := createOrgFileInOrg(ctx, entity)
	if err != nil {
		return result, err
	}

	return result, nil
}

func createOrgFileInOrg(ctx context.Context, entity entity.CreateOrgFileEntity) (int, error) {

	payload, _ := json.Marshal(entity)

	url := "http://" + serverUrl + "/org/server/v1/org/file"

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
