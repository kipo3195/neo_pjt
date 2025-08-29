package client

import (
	"context"
	"fmt"
	"org/internal/consts"
	"org/internal/domains/org/dto/client/requestDTO"
	"org/internal/domains/org/entities"
	"org/internal/domains/org/models"
	repositories "org/internal/domains/org/repositories/client"
	"org/internal/infra/storage"
	"strings"
)

type orgUsecase struct {
	repository repositories.OrgRepository
	orgStorage storage.OrgFileStorage
}

type OrgUsecase interface {
	GetOrgHash(ctx context.Context, req requestDTO.GetOrgHashRequest) (map[string]any, error)
	GetOrgData(ctx context.Context, req requestDTO.GetOrgDataRequest) (string, interface{}, error)
}

func NewOrgUsecase(repository repositories.OrgRepository, orgStorage storage.OrgFileStorage) OrgUsecase {
	return &orgUsecase{
		repository: repository,
		orgStorage: orgStorage,
	}
}

func (r *orgUsecase) GetOrgHash(ctx context.Context, req requestDTO.GetOrgHashRequest) (map[string]any, error) {

	orgMap := make(map[string]any)

	for i := 0; i < len(req.OrgHash); i++ {
		parts := strings.Split(req.OrgHash[i], "_")
		if len(parts) == 2 {

			fileFlag, eventFlag, err := r.repository.CheckOrgHash(ctx, parts[0], parts[1])

			if err != nil {
				return nil, err
			} else if fileFlag {
				// 파일로 받아야함.
				orgMap[req.OrgHash[i]] = "file"
			} else if eventFlag {
				// 이벤트로 받아야함.
				orgMap[req.OrgHash[i]] = "event"
			} else {
				// 최신 버전.
				orgMap[req.OrgHash[i]] = "latest"
			}
		} else {
			fmt.Printf("GetOrgs org : %s is invalid !", req.OrgHash[i])
			continue
		}
	}

	return orgMap, nil
}

func (r *orgUsecase) GetOrgData(ctx context.Context, req requestDTO.GetOrgDataRequest) (string, interface{}, error) {

	if req.Type == consts.FILE {
		version, err := r.repository.GetOrgLatestVersion(ctx, req.OrgCode)
		if err != nil {
			return "", nil, err
		}

		data, err := r.orgStorage.GetOrgFile(req.OrgCode)

		//filePath := "./storage/" + req.OrgCode + "/org_files/" + version // 전달할 파일 경로
		// 파일을 메모리에 가지고 있도록 수정 할 것.
		// fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일을 찾을 수 없음 %s \n", req.OrgCode)
			return "", nil, err
		}
		return version, data, nil

	} else if req.Type == consts.EVENT {

		events, err := r.repository.GetOrgDiffEvent(ctx, req.OrgCode, req.OrgHash)
		if err != nil {
			return "", nil, err
		}
		return "", toEventEntity(events), nil

	} else {
		// 명확하지 않은 타입으로 요청함.
		return "", nil, fmt.Errorf("invalid request type")
	}

}

func toEventEntity(events []models.OrgEvent) []entities.OrgEventEntity {

	var eventList []entities.OrgEventEntity

	for _, event := range events {
		entity := entities.OrgEventEntity{
			OrgCode:    event.OrgCode,
			EventType:  event.EventType,
			Id:         event.Id,
			Kind:       event.Kind,
			UpdateHash: event.UpdateHash,
		}
		eventList = append(eventList, entity)
	}
	return eventList
}
