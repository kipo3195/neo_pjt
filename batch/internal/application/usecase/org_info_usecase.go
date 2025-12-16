package usecase

import (
	"archive/zip"
	"batch/internal/application/util"
	"batch/internal/domain/orgInfo/entity"
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/storage"
	"bytes"
	"context"
	"encoding/json"
	"log"
)

type orgInfoUsecase struct {
	orgInfoStorage storage.OrgInfoStorage
	repo           repository.OrgInfoRepository
	apiRepo        repository.OrgInfoApiRepository
}

type OrgInfoUsecase interface {
	SendOrgInfoToOrg(ctx context.Context, org string) error
	PutSnapShotJson(ctx context.Context, org string, orgInfo []entity.OrgInfoEntity, fileName string) error
}

func NewOrgInfoUsecase(repo repository.OrgInfoRepository, apiRepo repository.OrgInfoApiRepository, orgInfoStorage storage.OrgInfoStorage) OrgInfoUsecase {

	return &orgInfoUsecase{
		orgInfoStorage: orgInfoStorage,
		repo:           repo,
		apiRepo:        apiRepo,
	}
}

func (r *orgInfoUsecase) SendOrgInfoToOrg(ctx context.Context, org string) error {

	// 현재 DB 조회 - 현재 조직도 json 파일 생성 zip 파일 생성 - 이걸 batch 서비스가 해야되는지 고민..

	// 현재 DB 조회 - zip 파일 생성
	orgInfo, err := r.repo.GetOrgInfo(ctx, org)

	if err != nil {
		return err
	}

	// 파일 명 생성
	fileName := util.GetNow() + ".json"
	log.Printf("[SendOrgInfoToOrg] org %s file name: %s\n", org, fileName)

	// 현재 DB 조회 데이터 json 생성 스냅샷 저장
	err = r.PutSnapShotJson(ctx, org, orgInfo, fileName)
	if err != nil {
		return err
	}

	// org 서비스 전송용
	orgInfoJson, err := json.MarshalIndent(orgInfo, "", "  ")
	if err != nil {
		return err
	}

	// json -> ZIP 파일 생성
	zipData, err := buildZipInMemory(fileName, orgInfoJson)
	if err != nil {
		return err
	}

	// 파일 전송
	err = r.apiRepo.SendJsonToOrg(ctx, fileName, zipData, org)
	if err != nil {
		return err
	}

	return nil
}

func parseOrgTree(orgTree []entity.OrgInfoEntity) *entity.OrgEntity {

	if orgTree == nil {
		log.Println("조회된 조직도 정보가 없음. ")
		return nil
	}

	// 최상위 구조
	var rootOrgInfos []entity.OrgInfo
	var flatList []entity.OrgInfo // 트리 구성용 전체 flat 리스트

	for _, org := range orgTree {
		// 이름 다국어 처리
		name := entity.NameEntity{
			Def: org.KoLang, // 수정 필요.
			Ko:  org.KoLang,
		}

		info := entity.OrgInfo{
			DeptCode:       org.DeptCode,
			ParentDeptCode: org.ParentDeptCode,
			Name:           name,
			Kind:           org.Kind,
			UserHash:       org.UserHash,
			UserId:         org.UserId,
			Header:         "",
			Description:    "",
		}

		if org.ParentDeptCode == "root" {
			rootOrgInfos = append(rootOrgInfos, info)
		}
		flatList = append(flatList, info)
	}

	// 트리 구조로 변환
	orgTreeInfo := buildOrgTree(flatList, "root")

	return &entity.OrgEntity{
		RootDept: rootOrgInfos,
		OrgTree:  orgTreeInfo,
	}
}

func buildOrgTree(flatList []entity.OrgInfo, parentCode string) []entity.OrgTreeInfo {
	var tree []entity.OrgTreeInfo

	for _, org := range flatList {
		if org.ParentDeptCode == parentCode {
			// 재귀적으로 하위 부서를 구성

			sub := buildOrgTree(flatList, org.DeptCode)

			// 사실 이렇게 구분해서 init하지 않아도 entity.OrgTreeInfo 내부에서 omitempty처리하면 response시 보이지 않음.
			if org.Kind == "0" {
				// 부서
				tree = append(tree, entity.OrgTreeInfo{
					DeptCode:       org.DeptCode,
					ParentDeptCode: org.ParentDeptCode,
					Name:           org.Name,
					SubDept:        sub,
					Kind:           org.Kind,
					UserHash:       org.UserHash,
					Description:    org.Description,
				})

				// 사용자
			} else if org.Kind == "1" {
				tree = append(tree, entity.OrgTreeInfo{
					ParentDeptCode: org.ParentDeptCode,
					Name:           org.Name,
					SubDept:        sub,
					Kind:           org.Kind,
					UserHash:       org.UserHash,
					UserId:         org.UserId,
					Header:         org.Header,
				})
			}

		}
	}

	return tree
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

func (r *orgInfoUsecase) PutSnapShotJson(ctx context.Context, org string, orgInfo []entity.OrgInfoEntity, fileName string) error {

	// 스냅샷 json 생성 (DB 저장용)
	orgInfoSnapShot := parseOrgTree(orgInfo)
	snapShot, err := json.MarshalIndent(orgInfoSnapShot, "", "  ")
	if err != nil {
		return err
	}

	// DB 백업
	err = r.repo.PutOrgInfoJson(ctx, org, fileName, string(snapShot))
	if err != nil {
		return err
	}
	return nil
}
