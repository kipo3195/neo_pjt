package usecases

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"org/consts"
	svDto "org/dto/server"
	"org/entities"
	"org/repositories"
	"os"
)

type serverUsecase struct {
	repo    repositories.ServerRepository
	orgRepo repositories.OrgRepository
	// 관심사 (책임)의 분리 측면에서 봤을때 GetOrg는 OrgRepo에서 처리. 동일한 두개의 로직을 만들지 않음.
}

func NewServerUsecase(repo repositories.ServerRepository) ServerUsecase {
	return &serverUsecase{
		repo: repo,
	}
}

type ServerUsecase interface {
	ServerCreateDept(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error)
	ServerDeleteDept(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error)

	ServerCreateDeptUser(ctx context.Context, req svDto.SvCreateDeptUserRequest) (interface{}, error)
	ServerDeleteDeptUser(ctx context.Context, req svDto.SvDeleteDeptUserRequest) (interface{}, error)

	ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error)
}

func (r *serverUsecase) ServerCreateDept(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error) {
	return r.repo.PutDept(ctx, toCreateDepartmentDto(req))
}

func toCreateDepartmentDto(req svDto.SvCreateDeptRequest) entities.CreateDeptEntity {

	return entities.CreateDeptEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		KrLang:         req.KrLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		CnLang:         req.CnLang,
	}
}

func (r *serverUsecase) ServerDeleteDeptUser(ctx context.Context, req svDto.SvDeleteDeptUserRequest) (interface{}, error) {
	return r.repo.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req svDto.SvDeleteDeptUserRequest) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *serverUsecase) ServerCreateDeptUser(ctx context.Context, req svDto.SvCreateDeptUserRequest) (interface{}, error) {

	updateHash := makeUpdateHash()
	fmt.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repo.PutDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req svDto.SvCreateDeptUserRequest, updateHash string) entities.CreateDeptUserEntity {

	return entities.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		DeptOrg:              req.DeptOrg,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
		UpdateHash:           updateHash,
	}
}

func (r *serverUsecase) ServerDeleteDept(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req svDto.SvDeleteDeptRequest) entities.DeleteDeptEntity {

	return entities.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *serverUsecase) ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error) {

	for i := 0; i < len(req.OrgCode); i++ {

		orgTree, err := r.repo.GetOrg(ctx, r.toGetOrgEntity(req.OrgCode[i]))

		if err != nil {
			fmt.Printf("ServerCreateOrgFile org : %s is invalid ! \n", req.OrgCode[i])
			continue
		}

		var hash = getNow()

		// 파일 명 생성.
		fileName := hash
		fmt.Printf("org %s file name : %s ", req.OrgCode[i], fileName)

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
		var zipPath = "./storage/" + req.OrgCode[i] + "/org_files/" + fileName

		// 1. 디렉터리 없으면 생성 (디렉터리 경로만 던져야함.) -> 실행시 마운트 필요
		err = ensureDir("./storage/" + req.OrgCode[i] + "/org_files/")
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
		fmt.Println("ZIP 파일 생성 ok path :", zipPath)

		// 3. ZIP writer 생성
		fmt.Println("ZIP writer 생성")
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()
		fmt.Println("ZIP writer 생성 ok")

		// 4. ZIP 내 파일 생성
		fmt.Println("ZIP 내 파일 생성")
		fileWriter, err := zipWriter.Create(fileName)
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

		result, err := r.repo.PutOrgEventHash(ctx, req.OrgCode[i], hash)
		if err != nil {
			return nil, fmt.Errorf("insert error %w", err)
		}
		if result {
			fmt.Println("DB saved ok org :", req.OrgCode[i])
		}
	}
	return consts.SUCCESS, nil
}

// 경로가 없으면 생성
func ensureDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

// 클라이언트, 서버에서 요청하여 현재의 ORG를 가져오기 위해 각자의 타입에 맞춰 entity 생성하는 코드 (오버라이드를 지원하지 않기 때문에 사용함)
func (r *serverUsecase) toGetOrgEntity(orgCode string) entities.GetOrgEntity {
	return entities.GetOrgEntity{
		OrgCode: orgCode,
	}
}
