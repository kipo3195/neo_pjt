package server

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"org/entities"
	"org/internal/domains/org/dto/server/requestDTO"
	repositories "org/internal/domains/org/repositories/server"
	"org/internal/infra/storage"
	"org/internal/utils"
	"org/models"
	"org/pkg/consts"
)

type orgUsecase struct {
	repository repositories.OrgRepository
	orgStorage storage.OrgFileStorage
}

type OrgUsecase interface {
	ServerCreateOrgFile(ctx context.Context, req requestDTO.CreateOrgFileRequest) (interface{}, error)
}

func NewOrgUsecase(repository repositories.OrgRepository, orgStorage storage.OrgFileStorage) OrgUsecase {
	return &orgUsecase{
		repository: repository,
		orgStorage: orgStorage,
	}
}

func (r *orgUsecase) ServerCreateOrgFile(ctx context.Context, req requestDTO.CreateOrgFileRequest) (interface{}, error) {

	for i := 0; i < len(req.OrgCode); i++ {

		org := req.OrgCode[i]

		orgTree, err := r.repository.GetOrg(ctx, org)
		if err != nil {
			fmt.Printf("Invalid org: %s\n", req.OrgCode[i])
			continue
		}

		// 저장시간 생성 = 파일 명
		fileName := utils.GetNow()
		fmt.Printf("org %s file name: %s\n", org, fileName)

		orgEntity := parseOrgTree(orgTree)
		orgJson, err := json.MarshalIndent(orgEntity, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal error: %w", err)
		}

		// 메모리에서 ZIP 생성
		zipData, err := buildZipInMemory(fileName, orgJson)
		if err != nil {
			return nil, fmt.Errorf("zip build error: %w", err)
		}

		// 메모리 저장소에 저장
		if err := r.orgStorage.SaveOrgFile(org, zipData); err != nil {
			return nil, fmt.Errorf("memory save error: %w", err)
		}

		// 점검
		data, err := r.orgStorage.GetOrgFile(org)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("저장된 org 파일의 사이즈 :", len(data))

		// DB 저장
		if ok, err := r.repository.PutOrgEventHash(ctx, org, fileName); err != nil {
			return nil, fmt.Errorf("db save error: %w", err)
		} else if ok {
			log.Println("DB saved ok org:", org)
		}
	}

	return consts.SUCCESS, nil
}

func buildZipInMemory(fileName string, content []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	writer, err := zipWriter.Create(fileName)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(content); err != nil {
		return nil, err
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseOrgTree(orgTree []models.WorksOrg) *entities.OrgEntity {

	if orgTree == nil {
		log.Println("조회된 조직도 정보가 없음. ")
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
