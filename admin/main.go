package main

import (
	"admin/internal/config"
	"admin/internal/domains/orgDeptUsers"
	"admin/internal/domains/orgDepts"
	"admin/internal/domains/orgFile"
	"admin/internal/domains/skinImg"
	"admin/internal/router"
	"log"
	"net/http"
)

func main() {

	log.Println("Admin service is running on :8089")
	server := InitServer()
	log.Fatal(server.ListenAndServe())
}

func InitServer() *http.Server {

	sfg := config.NewServerConfig()
	db := config.ConnectDatabase(sfg)

	r, baseGroup := router.SetDefaultRoutes("admin")

	// 스킨 이미지
	skinImgHandlers := skinImg.InitModules(db)
	router.SetSkinRoutes(baseGroup, skinImgHandlers)

	// 조직도 파일
	orgFileHandlers := orgFile.InitModules(db)
	router.SetOrgFilesRoutes(baseGroup, orgFileHandlers)

	// 부서에 사용자 추가
	orgDeptUsersHandlers := orgDeptUsers.InitModules(db)
	router.SetOrgDeptUsersRoutes(baseGroup, orgDeptUsersHandlers)

	// 부서 추가
	orgDeptsHandlers := orgDepts.InitModules(db)
	router.SetOrgDeptsRoutes(baseGroup, orgDeptsHandlers)

	return &http.Server{
		Addr:    ":8089",
		Handler: r,
	}

}
