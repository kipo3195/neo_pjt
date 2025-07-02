package usecases

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"org/consts"
	adminDto "org/dto/server/admin"
	"org/entities"
	storage "org/infra/storage"
	"org/repositories"
)

type serverUsecase struct {
	repo           repositories.ServerRepository
	orgFileStorage storage.OrgFileStorage
}

func NewServerUsecase(repo repositories.ServerRepository, orgFileStorage storage.OrgFileStorage) ServerUsecase {
	return &serverUsecase{
		repo:           repo,
		orgFileStorage: orgFileStorage,
	}
}

type ServerUsecase interface {
	ServerCreateDept(ctx context.Context, req adminDto.CreateDeptRequest) (interface{}, error)
	ServerDeleteDept(ctx context.Context, req adminDto.DeleteDeptRequest) (interface{}, error)

	ServerCreateDeptUser(ctx context.Context, req adminDto.CreateDeptUserRequest) (interface{}, error)
	ServerDeleteDeptUser(ctx context.Context, req adminDto.DeleteDeptUserRequest) (interface{}, error)

	ServerCreateOrgFile(ctx context.Context, req adminDto.CreateOrgFileRequest) (interface{}, error)
}

func (r *serverUsecase) ServerCreateDept(ctx context.Context, req adminDto.CreateDeptRequest) (interface{}, error) {
	return r.repo.PutDept(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(req adminDto.CreateDeptRequest) entities.CreateDeptEntity {

	return entities.CreateDeptEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		KoLang:         req.KoLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		ZhLang:         req.ZhLang,
		RuLang:         req.RuLang,
		ViLang:         req.ViLang,
		Header:         req.Header,
	}
}

func (r *serverUsecase) ServerDeleteDeptUser(ctx context.Context, req adminDto.DeleteDeptUserRequest) (interface{}, error) {
	return r.repo.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req adminDto.DeleteDeptUserRequest) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *serverUsecase) ServerCreateDeptUser(ctx context.Context, req adminDto.CreateDeptUserRequest) (interface{}, error) {

	updateHash := makeUpdateHash()
	fmt.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repo.PutDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req adminDto.CreateDeptUserRequest, updateHash string) entities.CreateDeptUserEntity {

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

func (r *serverUsecase) ServerDeleteDept(ctx context.Context, req adminDto.DeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req adminDto.DeleteDeptRequest) entities.DeleteDeptEntity {

	return entities.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *serverUsecase) ServerCreateOrgFile(ctx context.Context, req adminDto.CreateOrgFileRequest) (interface{}, error) {

	for i := 0; i < len(req.OrgCode); i++ {

		org := req.OrgCode[i]

		orgTree, err := r.repo.GetOrg(ctx, org)
		if err != nil {
			fmt.Printf("Invalid org: %s\n", req.OrgCode[i])
			continue
		}

		// 저장시간 생성 = 파일 명
		fileName := getNow()
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
		if err := r.orgFileStorage.SaveOrgFile(org, zipData); err != nil {
			return nil, fmt.Errorf("memory save error: %w", err)
		}

		// 점검
		data, err := r.orgFileStorage.GetOrgFile(org)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("저장된 org 파일의 사이즈 :", len(data))

		// DB 저장
		if ok, err := r.repo.PutOrgEventHash(ctx, org, fileName); err != nil {
			return nil, fmt.Errorf("db save error: %w", err)
		} else if ok {
			fmt.Println("DB saved ok org:", org)
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
