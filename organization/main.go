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
	departmentHandler := di.InitDepartmentHandler(db)
	router.SetDepartmentRoutes(baseGroup, departmentHandler.Handler)

	orgHandler := di.InitOrgHandler(db, orgStorage)
	router.SetOrgRoute(baseGroup, orgHandler.Handler)

	userHandler := di.InitUserHandler(db)
	router.SetUserRoute(baseGroup, userHandler.Handler)
	// ---- Orchestrator Init -----

	//router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: r,
	}
}
