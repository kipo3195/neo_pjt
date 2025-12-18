package usecase

import (
	"batch/internal/application/util"
	"batch/internal/domain/userDetail/repository"
	"context"
	"encoding/json"
	"log"
)

type userDetailUsecase struct {
	repo    repository.UserDetailRepository
	apiRepo repository.UserDetailApiRepository
}

type UserDetailUsecase interface {
	SendUserDetailToUser(ctx context.Context, org string) error
}

func NewUserDetailUsecase(repo repository.UserDetailRepository, apiRepo repository.UserDetailApiRepository) UserDetailUsecase {
	return &userDetailUsecase{
		repo:    repo,
		apiRepo: apiRepo,
	}
}

func (r *userDetailUsecase) SendUserDetailToUser(ctx context.Context, org string) error {

	// 현재 데이터 조회
	userDetail, err := r.repo.GetUserDetail(ctx, org)

	if err != nil {
		return err
	}

	// 파일 명 생성
	fileName := util.GetNow() + ".json"
	log.Printf("[SendUserDetailToUser] org %s file name: %s\n", org, fileName)

	// 현재 DB 조회 데이터 json 생성 스냅샷 저장
	// err = r.PutSnapShotJson(ctx, org, orgInfo, fileName)
	// if err != nil {
	// 	return err
	// }

	// org 서비스 전송용
	userDetailJson, err := json.MarshalIndent(userDetail, "", "  ")
	if err != nil {
		return err
	}

	// json -> ZIP 파일 생성
	zipData, err := buildZipInMemory(fileName, userDetailJson)
	if err != nil {
		return err
	}

	// 파일 전송
	err = r.apiRepo.SendJsonToUser(ctx, fileName, zipData, org)
	if err != nil {
		return err
	}

	return nil
}
