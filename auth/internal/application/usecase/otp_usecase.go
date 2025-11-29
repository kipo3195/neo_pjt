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

	// 기존 entity에서 Make 함수를 만들던 로직은 의존성 방향의 규칙을 깨는 일이므로 X
	devicePubKeyEntity := make([]entity.DevicePubKeyEntity, 0)

	for i := 0; i < len(input.DevicePubKey); i++ {
		e := entity.DevicePubKeyEntity{
			Kind: input.DevicePubKey[i].Kind,
			Key:  input.DevicePubKey[i].Key,
		}
		devicePubKeyEntity = append(devicePubKeyEntity, e)
	}

	entity := entity.MakeOtpKeyRegistEntity(input.Id, input.Uuid, devicePubKeyEntity)
	result, err := r.repo.OtpKeyRegistInMessage(ctx, entity)

	if err != nil {
		return output.OtpKeyRegistOutput{}, err
	}

	output := adapter.MakeOtpKeyRegistOutput(result)
	return output, nil
}
