package main

import (
	"log"
	"message/internal/delivery/router"
	"message/internal/di"
	"message/internal/infrastructure/config"
	"message/internal/infrastructure/migration"
	"message/internal/infrastructure/storage"
	"net/http"
)

func main() {
	log.Println("User service is running on :8084")
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

	// ---- Message Broker init ----
	mb := config.ConnectMessageBroker(sfg)

	// ---- Storage Init -----
	otpStorage := storage.NewOtpStorage()
	chatRoomStorage := storage.NewChatRoomStorage()
	// ---- Data Loader -----

	//dataLoader.Register(loader.NewDeviceTokenInfoLoader(db, deviceStorage))

	// ---- Router Init -----
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.
	router := router.NewMessageRouter("message", sfg.TokenConfig)

	// ---- Domain Handler Init -----

	noteModule := di.InitNoteModule(db, mb)
	router.SetNoteRoutes(noteModule.Handler)

	lineKeyModule := di.InitLineKeyModule(db)
	router.SetLineKeyRoutes(lineKeyModule.Handler)

	chatModule := di.InitChatModule(db, mb)
	router.SetChatRoutes(chatModule.Handler)

	otpModule := di.InitOtpModule(db, otpStorage)
	router.SetOtpRoutes(otpModule.Handler)

	chatRoomModule := di.InitChatRoomModule(db, chatRoomStorage)
	router.SetChatRoomRoutes(chatRoomModule.Handler)

	chatServiceModule := di.InitChatServiceModule(chatModule.Usecase, lineKeyModule.Usecase, chatRoomModule.Usecase)
	router.SetChatServiceRoutes(chatServiceModule)

	return &http.Server{
		Addr:    ":8083",
		Handler: router.GetEngine(),
	}
}
