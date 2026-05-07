package di

import (
	"net/http"
	"test/internal/adapter/http/router"
	"test/internal/infrastructure/config"
)

type AppContainer struct {
	Server *http.Server
}

func InitApp() (*AppContainer, error) {

	sfg := config.NewServerConfig()

	loginModule := InitLoginModule(sfg)

	router := router.NewTestRouter("test")
	router.SetLoginRoutes(loginModule.Handler)

	server := &http.Server{
		Addr:    ":8099",
		Handler: router.GetEngine(),
	}

	return &AppContainer{
		Server: server,
	}, nil

}
