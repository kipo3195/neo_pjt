package di

import (
	"batch/internal/adapter/scheduler"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/persistence/migration"
	"batch/internal/infrastructure/persistence/storage"
	"fmt"
	"log"
	"net/http"
)

type AppContainer struct {
	Server  *http.Server
	Cleanup func()
}

func InitApp() (*AppContainer, error) {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db, err := config.ConnectDatabase(sfg)
	if err != nil {
		return nil, fmt.Errorf("db connection failed:%w", err)
	}

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		err := migration.RunAll(db)
		if err != nil {
			return nil, fmt.Errorf("db auto migrate failed:%w", err)
		}
	}

	// ---- Storage Init -----
	orgInfoStorage := storage.NewOrgInfoStorage()
	userDetailStorage := storage.NewUserDetailStorage()

	// ---- Scheduler Init -----
	scheduler := scheduler.NewBatchScheduler(sfg)

	// ---- Domain Module Init -----
	orgInfoModule := InitOrgInfoModule(db, orgInfoStorage, sfg.Domain)
	userDetailModule := InitUserDetailModule(db, userDetailStorage, sfg.Domain)
	messageGrpcModule := InitMessageGrpcModule(db, sfg.ChatFileConfig)
	fileGrpcModule := InitFileGrpcModule(db, sfg.ChatFileConfig)

	extendDBConnectModule := InitExtendDBConnectModule(db)

	// ----- Service Orchestrator -----
	orgInfoBatchServiceModule := InitOrgInfoBatchServiceModule(orgInfoModule.Task, extendDBConnectModule.Task, sfg.OrgInfoBatchConfig)
	userDetailBatchServiceModule := InitUserDetailBatchserviceModule(userDetailModule.Task, extendDBConnectModule.Task, sfg.UserDetailBatchConfig)
	chatFileBatchServiceModule := InitChatFileBatchServiceModule(messageGrpcModule.Task, fileGrpcModule.Task)

	// ----- Scheduler Regist -----
	scheduler.RegistOrgInfoBatchService(orgInfoBatchServiceModule)
	scheduler.RegistUserDetailBatchService(userDetailBatchServiceModule)
	scheduler.RegistChatFileBatchService(chatFileBatchServiceModule)

	// ----- Scheduler Start -----
	scheduler.Start()

	// 서비스 자원정리
	cleanUp := func() {

		// 스케쥴러를 먼저 멈춰서 새로운 작업이 시작되지 않게함.
		if scheduler != nil {
			log.Println("Stop scheduler ...")
			scheduler.Stop()
		}

		// 만약 비동기 worker pool 처리시 여기에서 종료

		// 마지막에 DB 종료 (가장 안쪽 layer가 마지막에 정리되는 것이 바람직함. )
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			log.Println("Closing Database connection...")
			sqlDB.Close()
		}

	}

	server := &http.Server{
		Addr: ":8081",
	}

	return &AppContainer{
		Server:  server,
		Cleanup: cleanUp,
	}, nil
}
