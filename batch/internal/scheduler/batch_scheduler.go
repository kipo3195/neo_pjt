package scheduler

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
)

type batchScheduler struct {
	sfg *config.ServerConfig
}

type BatchScheduler interface {
	SetOrgInfoBatch(usecase usecase.OrgInfoBatchUsecase)
}

func NewBatchScheduler(sfg *config.ServerConfig) BatchScheduler {
	return &batchScheduler{
		sfg: sfg,
	}
}

func (r *batchScheduler) SetOrgInfoBatch(usecase usecase.OrgInfoBatchUsecase) {

	// sfg(서버 환경변수)에 있는 크론 실행 주기를 가지고 usecase 실행

}
