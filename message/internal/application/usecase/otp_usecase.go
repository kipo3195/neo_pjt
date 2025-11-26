package usecase

import (
	"context"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
	"message/internal/domain/otp/repository"

	"message/internal/infrastructure/storage"
)

type otpUsecase struct {
	repository repository.OtpRepository
	svChKey    string
	svNoKey    string
	otpStorage storage.OtpStorage
}

type OtpUsecase interface {
	OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output output.OtpKeyregistOutput, err error)
}

func NewOtpUsecase(repo repository.OtpRepository, storage storage.OtpStorage) OtpUsecase {
	return &otpUsecase{
		repository: repo,
		otpStorage: storage,
	}
}

func (u *otpUsecase) OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output output.OtpKeyregistOutput, err error) {

	// 키 생성 + 시간 생성,
	// DB 저장 (키 + 시간)
	// 메모리 저장 (storage)
	return output.OtpKeyregistOutput{
		ChkeyRegDate: "2024-01-01 00:00:00",
		NoKeyRegDate: "2024-01-01 00:00:00",
	}, err
}
