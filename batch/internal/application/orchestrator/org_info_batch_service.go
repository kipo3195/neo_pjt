package orchestrator

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
	"context"
	"log"
	"strconv"
	"time"
)

type OrgInfoBatchService struct {
	ExtendDBConnect usecase.ExtendDBConnectUsecase
	OrgInfo         usecase.OrgInfoUsecase
	serviceConfig   *config.OrgInfoBatchConfig
}

func NewOrgInfoServiceModule(orgInfo usecase.OrgInfoUsecase, extendDBConnection usecase.ExtendDBConnectUsecase, serviceConfig *config.OrgInfoBatchConfig) OrgInfoBatchService {

	return OrgInfoBatchService{
		OrgInfo:         orgInfo,
		ExtendDBConnect: extendDBConnection,
		serviceConfig:   serviceConfig,
	}
}

// 기존 handler 처리 로직의 트리거는 라우팅 된 handler에서 시작되었지만,
// 스케쥴러 처리는 scheduler에서 스케쥴링 등록, 시작 처리를 하고 Service에서 비즈니스 로직(usecase)의 트랜잭션을 갖도록 처리
// 그렇다면 기존의 orchestrator의 비즈니스 로직(여러 usecase의 묶음) 로직도 handler에서 실행하는 것이 아니라 서비스에서 해야되는거 아니냐?
// 맞음.. 그래서 수정예정
func (r *OrgInfoBatchService) Run(ctx context.Context) error {

	log.Printf("OrgInfoBatchService start. time : %s, extendDB : %s \n", time.Now().Format("2006-01-02 15:04:05"), strconv.FormatBool(r.serviceConfig.ExtendDBSyncFlag))

	// 외부 DB 커넥하여 데이터 조회 처리 로직
	if r.serviceConfig.ExtendDBSyncFlag {
		err := r.ExtendDBConnect.GetOrgInfo()
		if err != nil {
			log.Println("[OrgInfoBatchService] GetOrgInfo error : ", err)
		}
	}

	err := r.OrgInfo.SendOrgInfoToOrg(ctx, r.serviceConfig.Org)
	if err != nil {
		log.Println("[OrgInfoBatchService] SendOrgInfoToOrg error : ", err)
	}

	log.Println("OrgInfoBatchService end. time : ", time.Now().Format("2006-01-02 15:04:05"))

	return nil
}
