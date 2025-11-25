package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/domain/otp/entity"
)

func MakeOtpKeyRegistInput(id string, uuid string, chKey string, noKey string) input.OtpKeyRegistInput {
	return input.OtpKeyRegistInput{
		Id:    id,
		Uuid:  uuid,
		ChKey: chKey,
		NoKey: noKey,
	}
}

func MakeOtpKeyRegistOutput(entity entity.OtpKeyRegistResultEntity) output.OtpKeyRegistOutput {
	return output.OtpKeyRegistOutput{
		ChKeyRegDate: entity.ChKeyRegDate,
		NoKeyRegDate: entity.NoKeyRegDate,
	}
}
