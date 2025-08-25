package usecases

import (
	"log"
	"org/entities"
	storage "org/infra/storage"
	"org/internal/consts"
	"org/models"
	"org/repositories"
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
}

func NewOrgUsecase(repo repositories.OrgRepository, orgFileStorage storage.OrgFileStorage) OrgUsecase {
	return &orgUsecase{
		repo:           repo,
		orgFileStorage: orgFileStorage,
	}
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
