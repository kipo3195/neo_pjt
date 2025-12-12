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
	orgInfoBatchStorage := storage.NewOrgInfoBatchStorage()

	// ---- Data Loader -----

	// ---- Scheduler Init -----
	scheduler := scheduler.NewBatchScheduler(sfg)

	// ---- Domain Module Init -----
	orgInfoBatchModule := di.InitOrgInfoBatchModule(orgInfoBatchStorage)

	scheduler.SetOrgInfoBatch(orgInfoBatchModule.Usecase)

	return &http.Server{
		Addr: ":8081",
	}
}
