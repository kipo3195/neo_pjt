package main

import (
	"file/internal/delivery/handler"
	"file/internal/delivery/router"
	"file/internal/di"
	"file/internal/infrastructure/config"
	"file/internal/infrastructure/logger"
	"file/internal/infrastructure/migration"
	"file/internal/infrastructure/pb"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	// OCI SDK 필수 패키지
)

func main() {

	go func() {
		startGRPCServer()
	}()

	log.Println("file service is running on :8091")
	server := InitServer()
	log.Fatal(server.ListenAndServe())

}

// gRPC 서버 설정 및 실행 함수
func startGRPCServer() {
	lis, err := net.Listen("tcp", ":50001") // Message 서비스가 호출했던 그 포트
	if err != nil {
		log.Fatalf("gRPC 포트 열기 실패: %v", err)
	}

	s := grpc.NewServer()

	// 이 부분에서 실제 '서비스 핸들러'를 등록합니다.
	// pb는 protoc로 생성된 패키지입니다.
	// &grpcHandler{}는 아래 2번에서 만들 구현체입니다.
	pb.RegisterFileServiceServer(s, &handler.GrpcHandler{})

	log.Println("gRPC server is running on :50001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("gRPC 서버 실행 실패: %v", err)
	}
}

func InitServer() *http.Server {

	// ---- Server Config Init -----
	sfg := config.NewServerConfig()

	// ---- DB Connect -----
	db := config.ConnectDatabase(sfg)
	cacheClient := config.ConnectCacheDataBase(sfg)

	// ---- DB Migration -----
	if sfg.AutoMigrate {
		migration.RunAll(db)
	}

	// ---- LOGGER Init ----
	logger := logger.NewSlogLogger()

	// ---- Storage Init -----

	// ---- Data Loader -----

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewFileRouter("file", sfg.TokenConfig, logger)

	// ---- Domain Handler Init -----
	fileUrlModule := di.InitFileUrlModule(db, cacheClient, sfg.OracleStorageConfig, logger)
	router.SetFileUrlRoutes(fileUrlModule.Handler)

	return &http.Server{
		Addr:    ":8091",
		Handler: router.GetEngine(),
	}

}
