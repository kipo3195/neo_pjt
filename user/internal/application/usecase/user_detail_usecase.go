package usecase

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/consts"
	"user/internal/delivery/adapter"
	"user/internal/domain/userDetail/entity"
	"user/internal/domain/userDetail/repository"
	"user/internal/infrastructure/storage"
)

type userDetailUsecase struct {
	repository             repository.UserDetailRepository
	userInfoServiceStorage storage.UserInfoServiceStorage
}

type UserDetailUsecase interface {
	GetUserDetailInfo(ctx context.Context, input []input.GetUserDetailInfoInput) (output.GetUserDetailInfoOutput, error)
	RegisterUserDetailBatch(ctx context.Context, input input.RegistUserDetailBatchInput) error
}

func NewUserDatailUsecase(repository repository.UserDetailRepository, userInfoServiceStorage storage.UserInfoServiceStorage) UserDetailUsecase {
	return &userDetailUsecase{
		repository:             repository,
		userInfoServiceStorage: userInfoServiceStorage,
	}
}

func (u *userDetailUsecase) GetUserDetailInfo(ctx context.Context, input []input.GetUserDetailInfoInput) (output.GetUserDetailInfoOutput, error) {

	en := make([]entity.ReqUserEntity, 0)
	for _, i := range input {
		temp := entity.ReqUserEntity{
			UserHash:   i.UserHash,
			UpdateHash: i.UpdateHash,
		}
		en = append(en, temp)
	}
	// 기존 DB 조회 로직을 메모리 체크 로직으로 변경 TODO
	userInfos, err := u.repository.GetUserInfoDetailInfo(ctx, entity)

	if err != nil {
		return output.GetUserDetailInfoOutput{}, err
	}

	output := adapter.MakeGetUserDetailInfoOutput(userInfos)

	return output, nil
}

func (u *userDetailUsecase) RegisterUserDetailBatch(ctx context.Context, in input.RegistUserDetailBatchInput) error {

	en := entity.MakeRegistUserDetailBatchEntity(in.File, in.FileName, in.OrgCode)

	// zip 해제, json 구하기
	jsonBytes, err := unzipAndGetJSON(en.File)

	if err != nil {
		log.Println("[RegisterUserDetailBatch] unzipAndGetJSON error")
		return consts.ErrUnzipAndGetJSONError
	}

	// 2. JSON → Wrapper
	var orgInfo []entity.UserDetailBatchEntity
	if err := json.Unmarshal(jsonBytes, &orgInfo); err != nil {
		return consts.ErrInvalidUserDetailJSONError
	}

	log.Println("[RegisterUserDetailBatch] 연동 사용자 수  : ", len(orgInfo))

	err = u.repository.RegistUserDetail(ctx, orgInfo)

	if err != nil {
		return err
	}

	return nil
}

func unzipAndGetJSON(orgFile *[]byte) ([]byte, error) {

	// zip reader 생성
	zr, err := zip.NewReader(
		bytes.NewReader(*orgFile),
		int64(len(*orgFile)),
	)
	if err != nil {
		return nil, fmt.Errorf("zip reader error: %w", err)
	}

	// ZIP 내부 파일 순회
	for _, f := range zr.File {

		// json 파일만 추출
		if !strings.HasSuffix(f.Name, ".json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		// json 읽기
		jsonBytes, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		return jsonBytes, nil
	}

	return nil, fmt.Errorf("json file not found in zip")
}
