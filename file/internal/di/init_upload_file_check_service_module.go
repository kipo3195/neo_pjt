package di

import (
	"file/internal/adapter/rpc/grpcHandler"
	"file/internal/application/service"
	"file/internal/application/usecase"
)

func InitUploadFileCheckServiceModule(chatFile usecase.ChatFileUsecase, uploadFileCheck usecase.UploadFileCheckUsecase) *grpcHandler.UploadFileCheckServiceHandler {

	service := service.NewUploadFileCheckService(chatFile, uploadFileCheck)
	return grpcHandler.NewUploadFileCheckServiceHandler(service)

}
