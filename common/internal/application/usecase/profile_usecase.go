package usecase

import (
	"common/internal/application/usecase/input"
	"common/internal/application/usecase/output"
	"common/internal/application/util"
	"common/internal/consts"
	"common/internal/domain/profile/entity"
	"common/internal/domain/profile/repository"
	domainStorage "common/internal/domain/profile/storage"
	"common/internal/infrastructure/storage"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
)

type profileUsecase struct {
	repository          repository.ProfileRepository
	profileStorage      domainStorage.ProfileStorage
	profileCacheStorage storage.ProfileCacheStorage
}

type ProfileUsecase interface {
	ProfileImgUpload(ctx context.Context, in input.ProfileImgInput) error
	GetProfileImg(ctx context.Context, in input.GetProfileImgInput) (output.GetProfileImgOutput, error)
}

func NewProfileUsecase(repository repository.ProfileRepository, profileStorage domainStorage.ProfileStorage, profileCacheStorage storage.ProfileCacheStorage) ProfileUsecase {
	return profileUsecase{
		repository:          repository,
		profileStorage:      profileStorage,
		profileCacheStorage: profileCacheStorage,
	}
}

func (u profileUsecase) ProfileImgUpload(ctx context.Context, in input.ProfileImgInput) error {

	entity := entity.MakeProfileImgEntity(in.ProfileImg, in.ProfileImgSize, in.ProfileImgName, in.UserId)

	// 사이즈 체크
	sizeCheck := util.CheckProfileImgSize(entity.ProfileImgSize)
	if sizeCheck {
		return consts.ErrFileSizeExceeded
	}

	// 확장자 검증
	extentionCheck := util.ValidateImageFile(entity.ProfileImgName, *entity.ProfileImg)
	if !extentionCheck {
		return consts.ErrFileExtentionDetect
	}

	// 저장 파일명 생성 사용자 hash + 날짜
	profileImgHash := GenerateUserProfileHash(entity.UserId)
	log.Printf("[ProfileImgUpload] userId : %s, GenerateUserProfileHash : %s \n", entity.UserId, profileImgHash)

	// 파일 저장 처리 (저장 경로저장 필요시 _ 를 변수타입으로 변경해서 사용)
	saveFilePath, saveFileName, err := u.profileStorage.Upload(ctx, *entity.ProfileImg, profileImgHash)

	if err != nil {
		// 저장 에러 커스텀 에러 추가
		log.Printf("[ProfileImgUpload] %s file save error. \n", entity.UserId)
		return consts.ErrProfileImgSaveError
	}

	entity.ProfileImgSavedPath = saveFilePath
	entity.ProfileImgHash = profileImgHash
	entity.ProfileImgSavedName = saveFileName

	err = u.repository.PutUserProfileImgInfo(ctx, entity)
	if err != nil {
		log.Printf("[ProfileImgUpload] %s DB save error. \n", entity.UserId)
		// 파일 저장 삭제 처리 TODO
		return err
	}

	// id : 파일 명칭으로 저장
	u.profileCacheStorage.PutProfileName(entity.UserId, entity.ProfileImgSavedName)

	return nil
}

func GenerateUserProfileHash(userId string) string {
	date := time.Now().Format(consts.YYYYMMDDHHMSS)
	temp := userId + date
	hash := sha256.Sum256([]byte(temp))
	return hex.EncodeToString(hash[:])
}

func (u profileUsecase) GetProfileImg(ctx context.Context, in input.GetProfileImgInput) (output.GetProfileImgOutput, error) {

	entity := entity.MakeGetProfileImgEntity(in.UserId)

	profileName := u.profileCacheStorage.GetProfileName(entity.UserId)

	file, err := u.profileStorage.GetProfileUrl(ctx, profileName)

	if err != nil {
		return output.GetProfileImgOutput{}, err
	}

	output := output.MakeGetProfileImgOutput(file, profileName)

	return output, nil
}
