package usecases

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"org/consts"
	orgDto "org/dto/client/org"
	"org/entities"
	"org/infra/storage"
	"org/models"
	"org/repositories"
	"strings"
	"time"
)

type orgUsecase struct {
	repo           repositories.OrgRepository
	orgFileStorage storage.OrgFileStorage
}

func getNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}

type OrgUsecase interface {
	GetOrgHash(ctx context.Context, req orgDto.GetOrgHashRequest) (map[string]any, error)
	GetOrgData(ctx context.Context, req orgDto.GetOrgDataRequest) (string, interface{}, error)
}

func NewOrgUsecase(repo repositories.OrgRepository, orgFileStorage storage.OrgFileStorage) OrgUsecase {
	return &orgUsecase{
		repo:           repo,
		orgFileStorage: orgFileStorage,
	}
}

func (r *orgUsecase) GetOrgHash(ctx context.Context, req orgDto.GetOrgHashRequest) (map[string]any, error) {

	orgMap := make(map[string]any)

	for i := 0; i < len(req.OrgHash); i++ {
		parts := strings.Split(req.OrgHash[i], "_")
		if len(parts) == 2 {

			fileFlag, eventFlag, err := r.repo.CheckOrgHash(ctx, parts[0], parts[1])

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

func parseOrgTree(orgTree []models.WorksOrg) *entities.OrgEntity {

	if orgTree == nil {
		fmt.Println("조회된 조직도 정보가 없음. ")
		return nil
	}

	// 최상위 구조
	var rootOrgInfos []entities.OrgInfo
	var flatList []entities.OrgInfo // 트리 구성용 전체 flat 리스트

	for _, org := range orgTree {
		// 이름 다국어 처리
		name := entities.NameEntity{
			Def: org.KoLang, // 수정 필요.
			Ko:  org.KoLang,
			En:  org.EnLang,
			Jp:  org.JpLang,
			Zh:  org.ZhLang,
			Ru:  org.RuLang,
			Vi:  org.ViLang,
		}

		info := entities.OrgInfo{
			DeptCode:       org.DeptCode,
			ParentDeptCode: org.ParentDeptCode,
			Name:           name,
			Kind:           org.Kind,
			Id:             org.Id,
			Header:         org.Header,
		}

		if org.ParentDeptCode == "root" {
			rootOrgInfos = append(rootOrgInfos, info)
		}
		flatList = append(flatList, info)
	}

	// 트리 구조로 변환
	orgTreeInfos := buildOrgTree(flatList, "root")

	return &entities.OrgEntity{
		RootDept: rootOrgInfos,
		OrgTree:  orgTreeInfos,
	}
}

func buildOrgTree(flatList []entities.OrgInfo, parentCode string) []entities.OrgTreeInfos {
	var tree []entities.OrgTreeInfos

	for _, org := range flatList {
		if org.ParentDeptCode == parentCode {
			// 재귀적으로 하위 부서를 구성
			sub := buildOrgTree(flatList, org.DeptCode)

			tree = append(tree, entities.OrgTreeInfos{
				DeptCode:       org.DeptCode,
				ParentDeptCode: org.ParentDeptCode,
				Name:           org.Name,
				SubDept:        sub,
				Kind:           org.Kind,
				Id:             org.Id,
				Header:         org.Header,
			})
		}
	}

	return tree
}

func (r *orgUsecase) GetOrgData(ctx context.Context, req orgDto.GetOrgDataRequest) (string, interface{}, error) {

	if req.Type == consts.FILE {
		version, err := r.repo.GetOrgLatestVersion(ctx, req.OrgCode)
		if err != nil {
			return "", nil, err
		}

		data, err := r.orgFileStorage.GetOrgFile(req.OrgCode)

		//filePath := "./storage/" + req.OrgCode + "/org_files/" + version // 전달할 파일 경로
		// 파일을 메모리에 가지고 있도록 수정 할 것.
		// fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일을 찾을 수 없음 %s \n", req.OrgCode)
			return "", nil, err
		}
		return version, data, nil

	} else if req.Type == consts.EVENT {

		events, err := r.repo.GetOrgDiffEvent(ctx, req.OrgCode, req.OrgHash)
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

func makeUpdateHash() string {
	// 현재 시간 밀리초 문자열
	now := time.Now().UnixNano() // 나노초 단위 시간값

	// 16바이트 랜덤 바이트 생성
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 시간값을 바이트 배열로 변환 (int64 -> []byte)
	timeBytes := []byte(fmt.Sprintf("%d", now))

	// 시간 + 랜덤 바이트 합치기
	data := append(timeBytes, randomBytes...)

	// SHA-256 해시 생성
	hash := sha256.Sum256(data)

	// 16진수 인코딩해서 문자열 반환
	return hex.EncodeToString(hash[:])
}
