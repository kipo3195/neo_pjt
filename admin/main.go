package main

import (
	"admin/internal/delivery/router"
	"admin/internal/di"
	"admin/internal/infrastructure/config"
	"admin/internal/infrastructure/migration"
	"log"
	"net/http"
)

func main() {

	log.Println("Admin service is running on :8089")
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

	// ---- Data Loader -----

	// ---- Router Init -----
	router := router.NewAdminRouter("admin")
	// SetDefaultRoutes() 안에서 새로운 gin.Engine을 매번 생성하면 각기 다른 서버 인스턴스가 됩니다.
	// 이런 경우는 서버를 2개 띄우는 것과 같으므로 주의.

	// 스킨 이미지
	skinImgModule := di.InitSkinImgModule(db)
	router.SetSkinRoutes(skinImgModule.Handler)

	// 조직도 파일
	orgFileModule := di.InitOrgFileModule(db)
	router.SetOrgFileRoutes(orgFileModule.Handler)

	// 부서에 사용자 추가
	orgDeptUserModule := di.InitOrgDeptUserModule(db)
	router.SetOrgDeptUserRoutes(orgDeptUserModule.Handler)

	// 부서 추가
	orgDeptModule := di.InitOrgDeptModule(db)
	router.SetOrgDeptRoutes(orgDeptModule.Handler)

	userAuthRegisterModule := di.InitUserAuthRegisterModule()

	serviceUserModule := di.InitServiceUserModule(db)
	router.SetServiceUserRoutes(serviceUserModule.Handler)

	serviceUserAuthRegisterServiceModule := di.InitServiceUserAuthRegisterServiceModule(serviceUserModule.Usecase, userAuthRegisterModule.Usecase)
	router.SetServiceUserAuthRegisterServiceRoutes(serviceUserAuthRegisterServiceModule)

	return &http.Server{
		Addr:    ":8089",
		Handler: router.GetEngine(),
	}

}
