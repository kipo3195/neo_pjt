package main

import (
	"batch/internal/di"
	"batch/internal/infrastructure/config"
	"batch/internal/infrastructure/migration"
	"batch/internal/infrastructure/storage"
	"batch/internal/scheduler"
	"log"
	"net/http"
)

func main() {
	log.Println("Batch service is running on :8081")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- Storage Init -----
	orgInfoStorage := storage.NewOrgInfoStorage()
	userDetailStorage := storage.NewUserDetailStorage()

	// ---- Data Loader -----

	// ---- Scheduler Init -----
	scheduler := scheduler.NewBatchScheduler(sfg)

	// ---- Domain Module Init -----
	orgInfoModule := di.InitOrgInfoModule(db, orgInfoStorage, sfg.Domain)
	userDetailModule := di.InitUserDetailModule(db, userDetailStorage, sfg.Domain)

	extendDBConnectModule := di.InitExtendDBConnectModule(db)

	// ----- Service Orchestrator -----
	orgInfoBatchServiceModule := di.InitOrgInfoBatchServiceModule(orgInfoModule.Usecase, extendDBConnectModule.Usecase, sfg.OrgInfoBatchConfig)
	userDetailBatchServiceModule := di.InitUserDetailBatchserviceModule(userDetailModule.Usecase, extendDBConnectModule.Usecase, sfg.UserDetailBatchConfig)

	// ----- Scheduler Regist -----
	scheduler.RegistOrgInfoBatchService(orgInfoBatchServiceModule)
	scheduler.RegistUserDetailBatchService(userDetailBatchServiceModule)

	// ----- Scheduler Start -----
	scheduler.Start()

	return &http.Server{
		Addr: ":8081",
	}
}
