package di

import (
	"file/internal/adapter/rpc/grpcHandler"
	"file/internal/application/usecase"
	"file/internal/infrastructure/persistence/repository"

	"gorm.io/gorm"
)

type UploadFileCheckModule struct {
	UploadFilecheckGrpcHandler *grpcHandler.UploadFileCheckGrpcHandler
}

func InitUploadFileCheckModule(db *gorm.DB) UploadFileCheckModule {

	repository := repository.NewUploadFileCheckRepository(db)
	usecase := usecase.NewUploadFileCheckUsecase(repository)
	grpcHandler := grpcHandler.NewUploadFileCheckGrpcHandler(usecase)

	return UploadFileCheckModule{
		UploadFilecheckGrpcHandler: grpcHandler,
	}
}
