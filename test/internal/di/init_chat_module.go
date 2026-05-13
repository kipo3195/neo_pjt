package di

import (
	"test/internal/adapter/http/handler"
	"test/internal/application/usecase"
	"test/internal/infrastructure/config"
)

type ChatModule struct {
	Handler *handler.ChatHandler
	Usecase usecase.ChatUsecase
}

func InitChatModule(sfg *config.ServerConfig) *ChatModule {
	usecase := usecase.NewChatUsecase(sfg)
	handler := handler.NewChatHandler(usecase)

	return &ChatModule{
		Handler: handler,
		Usecase: usecase,
	}
}
