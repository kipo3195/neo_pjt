package scheduler

import (
	"batch/internal/application/service"
	"batch/internal/infrastructure/config"
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type batchScheduler struct {
	sfg *config.ServerConfig
	cr  *cron.Cron
}

type BatchScheduler interface {
	RegistOrgInfoBatchService(svc service.OrgInfoBatchService)
	RegistUserDetailBatchService(svc service.UserDetailBatchService)
	Start()
	Stop()
}

func NewBatchScheduler(sfg *config.ServerConfig) BatchScheduler {
	return &batchScheduler{
		sfg: sfg,
		cr:  cron.New(cron.WithSeconds()),
	}
}
func (r *batchScheduler) Stop() {
	r.cr.Stop()
}

// 등록과 실행의 책임을 완전히 분리
func (r *batchScheduler) Start() {
	r.cr.Start()
}

// 20251213 스케쥴러는 어떤 동작을 처리하는지 모르고 언제 스케쥴러를 등록하고 실행하는 책임만 갖도록..
func (r *batchScheduler) RegistOrgInfoBatchService(svc service.OrgInfoBatchService) {

	if r.sfg.OrgInfoBatchConfig.BatchFlag {

		// sfg(서버 환경변수)에 있는 크론 실행 주기를 가지고 usecase 실행
		cronExpr := r.sfg.OrgInfoBatchConfig.Cron // 예: "0 */5 * * *"

		log.Println("[SetOrgInfoBatch] cronExpr :", cronExpr)
		_, err := r.cr.AddFunc(cronExpr, func() {

			ctx := context.Background()
			//20251213 Orchestrator 각 비즈니스 로직에 대한 트랜잭션을 소유
			if err := svc.Run(ctx); err != nil {
				log.Println("[SetOrgInfoBatch] Org Sync Error:", err)
			}
		})

		if err != nil {
			log.Println("[SetOrgInfoBatch] cron config error : ", err)
		}

	} else {
		log.Println("[SetOrgInfoBatch] batch not used")
	}

}

func (r *batchScheduler) RegistUserDetailBatchService(svc service.UserDetailBatchService) {

	// sfg(서버 환경변수)에 있는 크론 실행 주기를 가지고 usecase 실행
	if r.sfg.UserDetailBatchConfig.BatchFlag {

		cronExpr := r.sfg.UserDetailBatchConfig.Cron // 예: "0 */5 * * *"

		log.Println("[SetUserDetailBatch] cronExpr :", cronExpr)
		_, err := r.cr.AddFunc(cronExpr, func() {

			ctx := context.Background()

			if err := svc.Run(ctx); err != nil {
				log.Println("[SetUserDetailBatch] Org Sync Error:", err)
			}
		})

		if err != nil {
			log.Println("[SetUserDetailBatch] cron config error : ", err)
		}
	} else {
		log.Println("[SetUserDetailBatch] batch not used")
	}
}
