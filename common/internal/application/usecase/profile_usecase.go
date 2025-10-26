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
	DeleteProfileImg(ctx context.Context, in input.DeleteProfileImgInput) error
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

	// 기존 프로필 삭제 로직 시작
	oldProfileName := u.profileCacheStorage.GetProfileName(entity.UserId)
	log.Println("[ProfileImgUpload] old Profile : ", oldProfileName)
	if oldProfileName != "" {
		// 이후 channel 로직으로 변경하기 (병렬처리)
		err = u.profileStorage.DeleteImg(ctx, oldProfileName)
		if err == nil {
			err = u.repository.DeleteUserProfileImgInfo(ctx, entity.UserId, oldProfileName)
			if err == consts.ErrProfileImgDBDeleteError || err == nil {
				u.profileCacheStorage.DeleteProfileName(entity.UserId, oldProfileName)
				log.Println("[ProfileImgUpload] old Profile delete success.")
			}
		}
	}
	// 기존 프로필 삭제 로직 끝

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

func (u profileUsecase) DeleteProfileImg(ctx context.Context, in input.DeleteProfileImgInput) error {

	entity := entity.MakeDeleteProfileImgEntity(in.UserId)
	profileName := u.profileCacheStorage.GetProfileName(entity.UserId)
	log.Println("[DeleteProfileImg] old Profile : ", profileName)
	if profileName == "" {
		// 프로필 이미지 등록되지 않은 사용자
		return consts.ErrProfileImgNotRegist
	}

	err := u.repository.DeleteUserProfileImgInfo(ctx, entity.UserId, profileName)

	if err == consts.ErrProfileImgDBDeleteError || err == nil {
		// 서버 경로 파일 삭제
		// 이후 channel 로직으로 변경하기 (병렬처리)
		err = u.profileStorage.DeleteImg(ctx, profileName)
		if err == nil {
			// DB에는 없지만 메모리에는 있는 case -> 삭제처리함.
			u.profileCacheStorage.DeleteProfileName(entity.UserId, profileName)
			log.Printf("[DeleteProfileImg] old Profile : %s delete success. \n", profileName)
		} else {
			// DB 삭제 실패, 파일 삭제 실패
			log.Println(err)
		}
	}
	return err
}
