package server

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"org/internal/domains/org/dto/server/requestDTO"
	repositories "org/internal/domains/org/repositories/server"
	"org/pkg/consts"
)

type orgUsecase struct {
	repository repositories.OrgRepository
}

type OrgUsecase interface {
	ServerCreateOrgFile(ctx context.Context, req requestDTO.CreateOrgFileRequest) (interface{}, error)
}

func NewOrgUsecase(repository repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{
		repository: repository,
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
		log.Println("저장된 org 파일의 사이즈 :", len(data))

		// DB 저장
		if ok, err := r.repo.PutOrgEventHash(ctx, org, fileName); err != nil {
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
