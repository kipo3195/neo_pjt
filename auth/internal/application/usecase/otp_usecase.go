package usecase

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/delivery/adapter"
	"auth/internal/domain/otp/entity"
	"auth/internal/domain/otp/repository"
	"context"
)

type otpUsecase struct {
	repo repository.OtpApiRepository
}

type OtpUsecase interface {
	OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output.OtpKeyRegistOutput, error)
}

func NewOtpUsecase(repo repository.OtpApiRepository) OtpUsecase {
	return &otpUsecase{
		repo: repo,
	}
}

func (r *otpUsecase) OtpKeyRegist(ctx context.Context, input input.OtpKeyRegistInput) (output.OtpKeyRegistOutput, error) {

	entity := entity.MakeOtpKeyRegistEntity(input.Id, input.Uuid, input.ChKey, input.NoKey)
	result, err := r.repo.OtpKeyRegistInMessage(ctx, entity)

	if err != nil {
		return output.OtpKeyRegistOutput{}, err
	}

	output := adapter.MakeOtpKeyRegistOutput(result)
	return output, nil
}
