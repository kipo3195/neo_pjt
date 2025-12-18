package orchestrator

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
	"context"
	"log"
	"strconv"
	"time"
)

type UserDetailBatchService struct {
	ExtendDBConnect usecase.ExtendDBConnectUsecase
	UserDetail      usecase.UserDetailUsecase
	serviceConfig   *config.BatchConfig
}

func NewUserDetailBatchServiceModule(userDetail usecase.UserDetailUsecase, extendDBConnection usecase.ExtendDBConnectUsecase, serviceConfig *config.BatchConfig) UserDetailBatchService {
	return UserDetailBatchService{
		UserDetail:      userDetail,
		ExtendDBConnect: extendDBConnection,
		serviceConfig:   serviceConfig,
	}
}

func (r *UserDetailBatchService) Run(ctx context.Context) error {

	log.Printf("UserDetailBatchService start. time : %s, extendDB : %s \n", time.Now().Format("2006-01-02 15:04:05"), strconv.FormatBool(r.serviceConfig.ExtendDBSyncFlag))

	// 외부 DB 커넥하여 데이터 조회 처리 로직
	if r.serviceConfig.ExtendDBSyncFlag {
		err := r.ExtendDBConnect.GetUserDetail()
		if err != nil {
			log.Println("[UserDetailBatchService] GetUserDetail error : ", err)
		}
	} else {
		log.Println("[UserDetailBatchService] extend db sync not used")
	}

	err := r.UserDetail.SendUserDetailToUser(ctx, r.serviceConfig.Org)
	if err != nil {
		log.Println("[UserDetailBatchService] SendUserDetailToUser error : ", err)
	}

	log.Println("UserDetailBatchService end. time : ", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
