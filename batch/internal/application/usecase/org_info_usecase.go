package usecase

import (
	"batch/internal/domain/orgInfo/repository"
	"batch/internal/infrastructure/storage"
	"log"
	"time"
)

type orgInfoUsecase struct {
	orgInfoStorage storage.OrgInfoStorage
	repo           repository.OrgInfoRepository
}

type OrgInfoUsecase interface {
	Run() error
}

func NewOrgInfoUsecase(repo repository.OrgInfoRepository, orgInfoStorage storage.OrgInfoStorage) OrgInfoUsecase {

	return &orgInfoUsecase{
		orgInfoStorage: orgInfoStorage,
		repo:           repo,
	}
}

func (r *orgInfoUsecase) Run() error {

	log.Println("[RUN] START time:", time.Now().Format("2006-01-02 15:04:05"))

	/* 트랜잭션 시작 */
	// 백업 옵션, json file, json DB
	// 현재 DB 조회 - 현재 조직도 json 파일 생성 zip 파일 생성

	// 외부 DB 조회 로직 -> TOBE 구조체 생성

	// ASIS DELETE + TOBE INSERT

	// 현재 DB 조회 - zip 파일 생성

	// server to server 전송 (to org)

	/* 트랜잭션 종료 */

	log.Println("[RUN] END time:", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
