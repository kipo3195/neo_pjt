package usecase

import (
	"context"
	"core/internal/application/usecase/input"
	"core/internal/domain/appValidation/entity"
	"core/internal/domain/appValidation/repository"
	"core/internal/infrastructure/storage"

	"core/internal/delivery/dto/appValidation"
	"fmt"
	"log"
)

type appValidationUsecase struct {
	repository        repository.AppValidationRepository
	serverInfoStorage storage.ServerInfoStorage
	apiRepository     repository.AppValidationAPIRepository
}

type AppValidationUsecase interface {
	CheckValidation(ctx context.Context, in input.AppValidationInput) (interface{}, error)
	GetWorksInfos(ctx context.Context, entity entity.ValidationEntity) (*appValidation.AppValidationResponseDTO, error)
}

func NewAppValidationUsecase(repository repository.AppValidationRepository, apiRepository repository.AppValidationAPIRepository, serverInfoStorage storage.ServerInfoStorage) AppValidationUsecase {

	return &appValidationUsecase{
		repository:        repository,
		apiRepository:     apiRepository,
		serverInfoStorage: serverInfoStorage,
	}
}

func (u *appValidationUsecase) CheckValidation(ctx context.Context, in input.AppValidationInput) (interface{}, error) { //output

	entity := entity.NewAppValidationEntity(in.Hash, in.Device, in.Uuid, in.WorksCode)

	err := u.repository.GetValidation(entity)

	if err != nil {
		return nil, err
	}

	result, err := u.GetWorksInfos(ctx, entity)

	return result, nil

}

// entitiy로 변경
func (u *appValidationUsecase) GetWorksInfos(ctx context.Context, entity entity.ValidationEntity) (*appValidation.AppValidationResponseDTO, error) {

	worksCode := entity.WorksCode

	// 메모리를 통한 조회
	worksCommonInfo := u.serverInfoStorage.GetWorksCommonInfo(worksCode)

	if worksCommonInfo == nil {
		fmt.Printf("[GetWorksCommonInfo] %s's worksCommonInfo is empty... check DB  \n", worksCode)
		// DB 조회를 통한 works 서버 정보 조회
		info, err := u.repository.GetWorksCommonInfo(worksCode)

		if err != nil {
			fmt.Printf("[GetWorksCommonInfo] worksCode %s's worksCommonInfo is DB empty... process end. \n", worksCode)
			return nil, err
		}
		u.serverInfoStorage.SaveWorksCommonInfo(worksCode, info)
		worksCommonInfo = info
	}

	// works의 domain/common API 호출 -> auth 호출 해서 jwt 발급, 저장, 결과 response.
	worksInfo, err := u.apiRepository.DeviceInit(ctx, entity)

	log.Println("worksInfo", worksInfo)

	if err != nil {
		log.Println("common service 호출시 에러 발생함.")
		return nil, err
	}

	return &appValidation.AppValidationResponseDTO{
		Body: appValidation.AppValidationResponseBody{
			WorksCommonInfo: worksCommonInfo,
			WorksInfo:       worksInfo,
		},
	}, nil

}
