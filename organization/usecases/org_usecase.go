package usecases

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"org/consts"
	clDto "org/dto/client"
	svDto "org/dto/server"
	"org/entities"
	"org/models"
	"org/repositories"
	"os"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
	GetOrg(ctx context.Context, req clDto.GetOrgRequest) (*entities.OrgEntity, error)

	ServerCreateDepartment(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error)
	ServerDeleteDepartment(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error)

	ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error)

	ToGetOrgEntity(req interface{}) entities.GetOrgEntity
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) GetOrg(ctx context.Context, req clDto.GetOrgRequest) (*entities.OrgEntity, error) {

	orgTree, err := r.repo.GetOrg(ctx, r.ToGetOrgEntity(req))
	if err != nil {
		return nil, err
	}

	orgEntity := parseOrgTree(orgTree)
	return orgEntity, nil
}

// 클라이언트, 서버에서 요청하여 현재의 ORG를 가져오기 위해 각자의 타입에 맞춰 entity 생성하는 코드 (오버라이드를 지원하지 않기 때문에 사용함)
func (r *orgUsecase) ToGetOrgEntity(req interface{}) entities.GetOrgEntity {
	switch v := req.(type) {
	case clDto.GetOrgRequest:
		return entities.GetOrgEntity{
			OrgCode: v.OrgCode,
		}
	case svDto.SvCreateOrgFileRequest:
		return entities.GetOrgEntity{
			OrgCode: v.OrgCode,
		}
	default:
		// 에러 처리하거나 빈 값 리턴
		return entities.GetOrgEntity{}
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

func (r *orgUsecase) ServerCreateDepartment(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error) {
	return r.repo.SaveDepartment(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(req svDto.SvCreateDeptRequest) entities.CreateDepartmentEntity {

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

func (r *orgUsecase) ServerDeleteDepartment(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDepartment(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req svDto.SvDeleteDeptRequest) entities.DeleteDepartmentEntity {

	return entities.DeleteDepartmentEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *orgUsecase) ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error) {

	orgTree, err := r.repo.GetOrg(ctx, r.ToGetOrgEntity(req))
	if err != nil {
		return nil, err
	}

	orgEntity := parseOrgTree(orgTree)
	// orgEntity를 JSON 등으로 직렬화 하고 내용을 ZIP 파일 내에 저장

	// 1. OrgEntity → JSON 직렬화
	fmt.Println("OrgEntity → JSON 직렬화")
	orgJson, err := json.MarshalIndent(orgEntity, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OrgEntity: %w", err)
	}
	fmt.Println("OrgEntity → JSON 직렬화 ok")

	// 경로에 파일명을 포함시켜야 함
	var zipPath = "./storage/org_files/org_entity.json"

	// 1. 디렉터리 없으면 생성 (디렉터리 경로만 던져야함.) -> 실행시 마운트 필요
	err = ensureDir("./storage/org_files/")
	if err != nil {
		return nil, fmt.Errorf("디렉터리 생성 실패: %w", err)
	}

	// 2. ZIP 파일 생성
	fmt.Println("ZIP 파일 생성")
	zipFile, err := os.Create(zipPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()
	fmt.Println("ZIP 파일 생성 ok")

	// 3. ZIP writer 생성
	fmt.Println("ZIP writer 생성")
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	fmt.Println("ZIP writer 생성 ok")

	// 4. ZIP 내 파일 생성
	fmt.Println("ZIP 내 파일 생성")
	fileWriter, err := zipWriter.Create("org_entity.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create file in zip: %w", err)
	}
	fmt.Println("ZIP 내 파일 생성 ok")

	// 5. write
	fmt.Println("Write")
	_, err = fileWriter.Write(orgJson)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to zip: %w", err)
	}
	fmt.Println("Write ok")
	return consts.SUCCESS, nil
}

// 경로가 없으면 생성
func ensureDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}
