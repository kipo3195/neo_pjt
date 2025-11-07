package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"time"
	"user/internal/application/usecase/input"
	"user/internal/application/usecase/output"
	"user/internal/application/util"
	"user/internal/consts"
	"user/internal/domain/profile/entity"
	"user/internal/domain/profile/repository"
	domainStorage "user/internal/domain/profile/storage"
	"user/internal/infrastructure/storage"
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
	RegistProfileMsg(ctx context.Context, in input.RegistProfileMsgInput) error
	GetProfileInfo(ctx context.Context, in input.GetProfileInfoInput) (output.GetProfileInfoOutput, error)
}

func NewProfileUsecase(repository repository.ProfileRepository, profileStorage domainStorage.ProfileStorage, profileCacheStorage storage.ProfileCacheStorage) ProfileUsecase {
	return &profileUsecase{
		repository:          repository,
		profileStorage:      profileStorage,
		profileCacheStorage: profileCacheStorage,
	}
}

func (u *profileUsecase) ProfileImgUpload(ctx context.Context, in input.ProfileImgInput) error {

	entity := entity.MakeProfileImgEntity(in.ProfileImg, in.ProfileImgSize, in.ProfileImgName, in.UserId, in.UserHash)

	// 메모리 관리의 key는 userHash로 처리 ]
	profileKey := entity.UserHash

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
	profileImgHash := GenerateUserProfileHash(profileKey)
	log.Printf("[ProfileImgUpload] userHash : %s, GenerateUserProfileHash : %s \n", profileKey, profileImgHash)

	// 파일 저장 처리 (저장 경로저장 필요시 _ 를 변수타입으로 변경해서 사용)
	saveFilePath, saveFileName, err := u.profileStorage.Upload(ctx, *entity.ProfileImg, profileImgHash)

	if err != nil {
		// 저장 에러 커스텀 에러 추가
		log.Printf("[ProfileImgUpload] %s file save error. \n", profileKey)
		return consts.ErrProfileImgSaveError
	}

	// 기존 프로필 삭제 로직 시작
	oldProfileName := u.profileCacheStorage.GetProfileName(profileKey)
	// DB 조회
	if oldProfileName != "" {
		log.Println("[ProfileImgUpload] 기존 프로필 이미지 삭제 프로세스 시작")
		err = u.repository.DeleteUserProfileImgInfo(ctx, profileKey, oldProfileName)
		// 메모리에는 있는데 DB에 없음 -> 정상적이지 않은 파일로 간주
		if err == consts.ErrProfileImgDBDeleteError {

			// 서버 경로 파일 삭제
			u.profileStorage.DeleteImg(ctx, oldProfileName)

			// 메모리 삭제
			u.profileCacheStorage.DeleteProfileName(profileKey, oldProfileName)

			// DB 처리 완료됨
		} else if err == nil {
			// 서버 경로 파일 삭제
			err = u.profileStorage.DeleteImg(ctx, oldProfileName)

			// 삭제 에러 발생함
			if err == consts.ErrProfileImgRemoveError {
				// DB roll back, 메모리 삭제 X
				log.Println("[ProfileImgUpload] 기존 프로필 이미지 삭제 불가.. 여기가 쌓이면 문제 생김")
				u.repository.RollbackDeleteUserProfileImgInfo(ctx, profileKey, oldProfileName)
			} else {
				// 정상적으로 삭제됨
				// 파일이 존재하지 않든, 삭제하다 실패하든 상관없이 메모리 삭제 처리
				u.profileCacheStorage.DeleteProfileName(profileKey, oldProfileName)
			}
		}
		log.Println("[ProfileImgUpload] 기존 프로필 이미지 삭제 프로세스 종료 userHash : ", profileKey)
	}

	// 기존 프로필 삭제 로직 끝

	entity.ProfileImgSavedPath = saveFilePath
	entity.ProfileImgHash = profileImgHash
	entity.ProfileImgSavedName = saveFileName

	err = u.repository.PutUserProfileImgInfo(ctx, entity)
	if err != nil {
		log.Printf("[ProfileImgUpload] %s DB save error. \n", profileKey)
		// 파일 저장 삭제 처리 TODO
		return err
	}

	// id : 파일 명칭으로 저장
	u.profileCacheStorage.PutProfileName(profileKey, entity.ProfileImgSavedName)

	return nil
}

func GenerateUserProfileHash(userId string) string {
	date := time.Now().Format(consts.YYYYMMDDHHMSS)
	temp := userId + date
	hash := sha256.Sum256([]byte(temp))
	return hex.EncodeToString(hash[:])
}

func (u *profileUsecase) GetProfileImg(ctx context.Context, in input.GetProfileImgInput) (output.GetProfileImgOutput, error) {

	entity := entity.MakeGetProfileImgEntity(in.UserId)

	profileName := u.profileCacheStorage.GetProfileName(entity.UserHash)

	if profileName == "" {
		// 메모리에 없음 = 실제 파일 없음.
		return output.GetProfileImgOutput{}, consts.ErrProfileImgNotRegist
	}

	file, err := u.profileStorage.GetProfileUrl(ctx, profileName)

	if err != nil {
		return output.GetProfileImgOutput{}, err
	}

	output := output.MakeGetProfileImgOutput(file, profileName)

	return output, nil
}

func (u *profileUsecase) DeleteProfileImg(ctx context.Context, in input.DeleteProfileImgInput) error {

	entity := entity.MakeDeleteProfileImgEntity(in.UserId)
	profileName := u.profileCacheStorage.GetProfileName(entity.UserId)
	log.Println("[DeleteProfileImg] old Profile : ", profileName)
	if profileName == "" {
		// 프로필 이미지 등록되지 않은 사용자
		return consts.ErrProfileImgNotRegist
	}

	// DB 조회
	err := u.repository.DeleteUserProfileImgInfo(ctx, entity.UserId, profileName)

	// 메모리에는 있는데 DB에 없음 -> 정상적이지 않은 파일로 간주
	if err == consts.ErrProfileImgDBDeleteError {
		// 서버 경로 파일 삭제
		err = u.profileStorage.DeleteImg(ctx, profileName)

		// 메모리 삭제
		u.profileCacheStorage.DeleteProfileName(entity.UserId, profileName)

		// DB 처리 완료됨
	} else if err == nil {
		// 서버 경로 파일 삭제
		err = u.profileStorage.DeleteImg(ctx, profileName)

		// 삭제 에러 발생함
		if err == consts.ErrProfileImgRemoveError {
			log.Println("[DeleteProfileImg] 기존 프로필 이미지 삭제 불가.. 여기가 쌓이면 문제 생김")
			// DB roll back, 메모리 삭제 X
			u.repository.RollbackDeleteUserProfileImgInfo(ctx, entity.UserId, profileName)
		} else {
			// 정상적으로 삭제됨
			// 파일이 존재하지 않든, 삭제하다 실패하든 상관없이 메모리 삭제 처리
			u.profileCacheStorage.DeleteProfileName(entity.UserId, profileName)
		}
	}

	return err
}

func (u *profileUsecase) RegistProfileMsg(ctx context.Context, in input.RegistProfileMsgInput) error {

	entity.MakePutProfileMsgEntity(in.UserId, in.Msg)

	return nil

}

func (u *profileUsecase) GetProfileInfo(ctx context.Context, in input.GetProfileInfoInput) (output.GetProfileInfoOutput, error) {

	entity := entity.MakeGetProfileInfoEntity(in.UserHashs)
	profileInfo, err := u.repository.GetProfileInfo(ctx, entity)

	if err != nil {
		return output.GetProfileInfoOutput{}, err
	}

	return output.GetProfileInfoOutput{
		ResultMap: profileInfo,
	}, nil
}
