package usecases

import (
	"context"
	clDto "org/dto/client"
	svDto "org/dto/server"
	"org/entities"
	"org/models"
	"org/repositories"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
	GetOrg(ctx context.Context, req clDto.GetOrgRequest) (*entities.OrgEntity, error)

	ServerCreateDepartment(ctx context.Context, req svDto.ServerCreateDeptRequest) (interface{}, error)
	ServerDeleteDepartment(ctx context.Context, req svDto.ServerDeleteDeptRequest) (interface{}, error)
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) GetOrg(ctx context.Context, req clDto.GetOrgRequest) (*entities.OrgEntity, error) {

	orgTree, err := r.repo.GetOrg(ctx, toGetOrgEntity(req))
	if err != nil {
		return nil, err
	}

	orgEntity := parseOrgTree(orgTree)
	return orgEntity, nil
}

func toGetOrgEntity(req clDto.GetOrgRequest) entities.GetOrgEntity {
	return entities.GetOrgEntity{
		OrgCode: req.OrgCode,
	}
}

func parseOrgTree(orgTree *[]models.WorksOrg) *entities.OrgEntity {

	// 최상위 구조
	var rootOrgInfos []entities.OrgInfo
	var flatList []entities.OrgInfo // 트리 구성용 전체 flat 리스트

	for _, org := range *orgTree {
		info := entities.OrgInfo{
			DeptCode:       org.DeptCode,
			ParentDeptCode: org.ParentDeptCode,
			KrLang:         org.KrLang,
			EnLang:         org.EnLang,
			JpLang:         org.JpLang,
			CnLang:         org.CnLang,
			DeptUpdateHash: org.DeptUpdateHash,
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
				KrLang:         org.KrLang,
				EnLang:         org.EnLang,
				JpLang:         org.JpLang,
				CnLang:         org.CnLang,
				DeptUpdateHash: org.DeptUpdateHash,
				SubDept:        sub,
			})
		}
	}

	return tree
}

func (r *orgUsecase) ServerCreateDepartment(ctx context.Context, req svDto.ServerCreateDeptRequest) (interface{}, error) {
	return r.repo.SaveDepartment(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(req svDto.ServerCreateDeptRequest) entities.CreateDepartmentEntity {

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

func (r *orgUsecase) ServerDeleteDepartment(ctx context.Context, req svDto.ServerDeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDepartment(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req svDto.ServerDeleteDeptRequest) entities.DeleteDepartmentEntity {

	return entities.DeleteDepartmentEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}
