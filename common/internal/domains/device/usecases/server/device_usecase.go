package server

import (
	"common/entities"
	"common/internal/domains/device/dto/server/requestDTO"
	"common/internal/domains/device/repositories/serverRepository"
	"common/pkg/consts"
	"context"
	"log"
)

type deviceUsecase struct {
	repository serverRepository.DeviceRepository
}

type DeviceUsecase interface {
	DeviceInit(ctx context.Context, body *requestDTO.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse)
}

func NewDeviceUsecase(repository serverRepository.DeviceRepository) DeviceUsecase {
	return &deviceUsecase{
		repository: repository,
	}
}

func (u *deviceUsecase) DeviceInit(ctx context.Context, body *requestDTO.DeviceInitRequest) (*entities.InitResult, *dto.ErrorResponse) {

	// DB 조회 connectInfo(접속 url)은 관리해야할 필요있음. 최초 이후에 클라이언트가 정보가 필요할때를 대비해서.
	connectInfo, err := u.repository.GetConnectInfo(body.WorksCode)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_102,
			Message: consts.E_102_MSG,
		}
	}

	// AUTH에 JWT 요청
	issuedAppToken, err := generateAppToken(body, connectInfo.ServerUrl) //  serverUrl은 이후에 .env 또는 k8s의 secrets에서 읽기
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	// 타임존, 언어, 앱 별 스킨 정보, 설정 정보 - GetConnectInfo와 합쳐서 트랜잭션 처리
	worksConfig, err := u.repository.GetWorksConfig(toWorksConfigEntity(body.WorksCode, body.Device), ctx)
	if err != nil {
		return nil, &dto.ErrorResponse{
			Code:    consts.E_500,
			Message: consts.E_500_MSG,
		}
	}

	log.Println("connectInfo:", connectInfo)
	log.Println("issuedAppToken:", issuedAppToken)
	log.Println("worksConfig ", worksConfig)

	return &entities.InitResult{
		ConnectInfo:    connectInfo,
		IssuedAppToken: issuedAppToken,
		WorksConfig:    worksConfig,
	}, nil
}

func toWorksConfigEntity(worksCode string, device string) entities.GetWorksConfig {

	return entities.GetWorksConfig{
		WorksCode: worksCode,
		Device:    device,
	}

}
