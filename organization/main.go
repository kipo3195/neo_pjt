package main

import (
	"log"
	"net/http"
	"org/internal/app"
	"org/internal/config"
	"org/internal/delivery/router"
	"org/internal/infrastructure/migration"
	"org/internal/infrastructure/storage"
)

func main() {
	log.Println("org service is running on :8088")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
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

	departmentHandler := app.InitDepartmentHandler(db)
	router.SetDepartmentRoutes(baseGroup, departmentHandler.Handler)

	orgHandler := app.InitOrgHandler(db, orgStorage)
	router.SetOrgRoute(baseGroup, orgHandler.Handler)

	userHandler := app.InitUserHandler(db)
	router.SetUserRoute(baseGroup, userHandler.Handler)

	// ---- Service Init -----

	//router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: r,
	}
}
