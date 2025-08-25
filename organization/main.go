package main

import (
	"log"
	"net/http"
	"org/internal/config"
	"org/internal/domains/department"
	"org/internal/domains/org"
	"org/internal/domains/user"
	"org/internal/infra/storage"
	"org/internal/router"
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

	// ---- Storage Init -----
	orgStorage := storage.NewOrgFileStorage() // 조직도 메모리 관리

	// ---- Data Loader -----

	// ---- Router Init -----

	r, baseGroup := router.SetDefaultRoutes("org")

	departmentHandler := department.InitModule(db)
	router.SetDepartmentRoutes(baseGroup, departmentHandler)

	orgHandler := org.InitModule(db, orgStorage)
	router.SetOrgRoute(baseGroup, orgHandler)

	userHandler := user.InitModule(db)
	router.SetUserRoute(baseGroup, userHandler)

	// ---- Service Init -----

	//router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: r,
	}
}
