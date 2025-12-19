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
	repository repository.UserDetailRepository
	storage    storage.UserInfoServiceStorage
}

type UserDetailUsecase interface {
	GetUserDetailInfo(ctx context.Context, input []input.GetUserDetailInfoInput) ([]output.UserDetailOutput, error)
	GetMyDetailInfo(ctx context.Context, in []string) ([]output.UserDetailOutput, error)
	RegisterUserDetailBatch(ctx context.Context, input input.RegistUserDetailBatchInput) error
}

func NewUserDatailUsecase(repository repository.UserDetailRepository, storage storage.UserInfoServiceStorage) UserDetailUsecase {
	return &userDetailUsecase{
		repository: repository,
		storage:    storage,
	}
}

func (u *userDetailUsecase) GetMyDetailInfo(ctx context.Context, in []string) ([]output.UserDetailOutput, error) {

	detailInfo := u.storage.GetUserDetailInfo(in)

	output := adapter.MakeGetUserDetailInfoOutput(detailInfo)

	return output, nil
}

func (u *userDetailUsecase) GetUserDetailInfo(ctx context.Context, input []input.GetUserDetailInfoInput) ([]output.UserDetailOutput, error) {

	// entity 생성
	en := make([]entity.ReqUserEntity, 0)
	for _, i := range input {
		temp := entity.ReqUserEntity{
			UserHash:   i.UserHash,
			DetailHash: i.DetailHash,
		}
		en = append(en, temp)
	}
	// 기존 DB 조회 로직을 메모리 체크 로직으로 변경 TODO
	//userInfos, err := u.repository.GetUserInfoDetailInfo(ctx, entity)

	// 이 방식이 루프 안에서 N번 호출하는 것보다 훨씬 빠릅니다.
	currentInfos, err := u.storage.GetUserInfoUpdateHash(ctx, en)
	if err != nil {
		return nil, err
	}

	// 요청한 사용자의 updateHash 데이터 비교
	targetUsers := make([]string, 0)
	for _, req := range en {
		current, exists := currentInfos[req.UserHash]

		// 변경된 것만 골라냄 (이전 질문에서 논의한 직접 비교)
		if !exists {
			// 없다 - DB 한번 더 체크?
			log.Printf("[GetUserDetailInfo] %s hash is not exist \n", req.UserHash)
		} else if current.DetailHash != req.DetailHash {
			// 있지만 다르다.
			targetUsers = append(targetUsers, req.UserHash)
		}
	}

	detailInfo := u.storage.GetUserDetailInfo(targetUsers)

	output := adapter.MakeGetUserDetailInfoOutput(detailInfo)

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
