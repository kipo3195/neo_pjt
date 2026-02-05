package di

import (
	"file/internal/adapter/http/handler"
	"file/internal/adapter/rpc/grpcHandler"
	"file/internal/application/usecase"
	"file/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type ChatFileModule struct {
	ChatFileHandler     *handler.ChatFileHandler
	ChatFileGrpcHandler *grpcHandler.ChatFileGrpcHandler
}

func InitChatFileModule(db *gorm.DB) *ChatFileModule {

	repository := repository.NewChatFileRepository(db)
	usecase := usecase.NewChatFileUsecase(repository)

	handler := handler.NewChatFileHandler(usecase)
	grpcHandler := grpcHandler.NewChatFileGrpcHandler(usecase)

	return &ChatFileModule{
		ChatFileHandler:     handler,
		ChatFileGrpcHandler: grpcHandler,
	}
}
