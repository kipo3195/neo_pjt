package main

import (
	"context"
	"log"
	"net/http"
	"org/config"
	"org/handlers"
	"org/infra/storage"
	"org/internal/domains/department"
	"org/internal/router"
	"org/repositories"
	"org/usecases"
	"time"
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
	orgFileStorage := storage.NewOrgFileStorage() // 조직도 메모리 관리

	// ---- Data Loader -----
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ---- Router Init -----

	r, baseGroup := router.SetDefaultRoutes("org")

	departmentHandler := department.InitModule(db)
	router.SetDepartmentRoutes(baseGroup, departmentHandler)

	// ---- Service Init -----

	orgRepo := repositories.NewOrgRepository(db)
	orgUsecase := usecases.NewOrgUsecase(orgRepo, orgFileStorage)
	orgHandler := handlers.NewOrgHandler(sfg, orgUsecase)

	userRepo := repositories.NewUserRepository(db)
	userUsecase := usecases.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(sfg, userUsecase)

	serverRepo := repositories.NewServerRepository(db)
	serverUsecase := usecases.NewServerUsecase(serverRepo, orgFileStorage)
	serverHandler := handlers.NewServerHandler(sfg, serverUsecase)

	//router := routes.SetupRoutes(handlers)

	return &http.Server{
		Addr:    ":8088",
		Handler: r,
	}
}
