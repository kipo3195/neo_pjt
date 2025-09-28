package main

import (
	"log"
	"net/http"
	"org/internal/delivery/router"
	"org/internal/di"
	"org/internal/infrastructure/config"
	"org/internal/infrastructure/migration"
	"org/internal/infrastructure/storage"
)

func main() {
	log.Println("org service is running on :8088")
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
	orgStorage := storage.NewOrgFileStorage() // 조직도 메모리 관리

	// ---- Data Loader -----

	// ---- Router Init -----
	r, baseGroup := router.SetDefaultRoutes("org")

	// ---- Domain Handler Init -----
	departmentModule := di.InitDepartmentModule(db)
	router.SetDepartmentRoutes(baseGroup, departmentModule.Handler, sfg.TokenConfig)

	orgModule := di.InitOrgModule(db, orgStorage)
	router.SetOrgRoute(baseGroup, orgModule.Handler, sfg.TokenConfig)

	userModule := di.InitUserModule(db)
	router.SetUserRoute(baseGroup, userModule.Handler, sfg.TokenConfig)
	// ---- Orchestrator Init -----

	//router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: r,
	}
}
