package di

import (
	"message/internal/application/usecase"
	"message/internal/delivery/handler"
	"message/internal/infrastructure/repository"
	"message/internal/infrastructure/storage"

	"gorm.io/gorm"
)

type OtpModule struct {
	Handler *handler.OtpHandler
}

func InitOtpModule(db *gorm.DB, storage storage.OtpStorage) *OtpModule {
	repository := repository.NewOtpApiRepository(db)
	// 테스트 용이므로 하드코딩, 실제 서비스 시 환경변수 등에서 주입 필요
	usecase := usecase.NewOtpUsecase(repository, storage, "your_service_chat_otp_key", "v1", "your_service_note_otp_key", "v1")
	handler := handler.NewOtpHandler(usecase)
	return &OtpModule{
		Handler: handler,
	}
}
