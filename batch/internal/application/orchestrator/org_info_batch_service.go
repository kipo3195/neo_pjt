package orchestrator

import (
	"batch/internal/application/usecase"
	"context"
	"log"
	"time"
)

type OrgInfoBatchService struct {
	ExtendDBConnect usecase.ExtendDBConnectUsecase
	OrgInfo         usecase.OrgInfoUsecase
}

func NewOrgInfoServiceModule(orgInfo usecase.OrgInfoUsecase, extendDBConnection usecase.ExtendDBConnectUsecase) OrgInfoBatchService {

	return OrgInfoBatchService{
		OrgInfo:         orgInfo,
		ExtendDBConnect: extendDBConnection,
	}
}

func (r *OrgInfoBatchService) Run(ctx context.Context) error {
	log.Println("OrgInfoBatchService start time : ", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}
