package service

import "file/internal/application/usecase"

type UploadFilecheckService struct {
	ChatFile        usecase.ChatFileUsecase
	UploadFileCheck usecase.UploadFileCheckUsecase
}

func NewUploadFileCheckService(chatFile usecase.ChatFileUsecase, uploadFileCheck usecase.UploadFileCheckUsecase) UploadFilecheckService {
	return UploadFilecheckService{
		ChatFile:        chatFile,
		UploadFileCheck: uploadFileCheck,
	}
}
