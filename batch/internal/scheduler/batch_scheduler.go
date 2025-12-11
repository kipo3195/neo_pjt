package scheduler

import (
	"batch/internal/application/usecase"
	"batch/internal/infrastructure/config"
	"log"

	"github.com/robfig/cron/v3"
)

type batchScheduler struct {
	sfg *config.ServerConfig
	cr  *cron.Cron
}

type BatchScheduler interface {
	SetOrgInfoBatch(usecase usecase.OrgInfoBatchUsecase)
}

func NewBatchScheduler(sfg *config.ServerConfig) BatchScheduler {
	return &batchScheduler{
		sfg: sfg,
		cr:  cron.New(cron.WithSeconds()),
	}
}

func (r *batchScheduler) SetOrgInfoBatch(usecase usecase.OrgInfoBatchUsecase) {

	// sfg(서버 환경변수)에 있는 크론 실행 주기를 가지고 usecase 실행
	cronExpr := r.sfg.OrgInfoBatchCron // 예: "0 */5 * * *"

	log.Println("[SetOrgInfoBatch] cronExpr :", cronExpr)
	_, err := r.cr.AddFunc(cronExpr, func() {
		log.Println("[Batch] Run Org Info Sync Batch")
		err := usecase.Run()
		if err != nil {
			log.Println("[Batch] Org Sync Error:", err)
		}
	})

	log.Println("[SetOrgInfoBatch] start ")
	r.cr.Start()

	if err != nil {
		log.Println("[SetOrgInfoBatch] cron config error : ", err)
	}

}
