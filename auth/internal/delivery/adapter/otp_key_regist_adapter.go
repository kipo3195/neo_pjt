package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
	"auth/internal/delivery/dto/deviceAuthService"
	"auth/internal/domain/otp/entity"
)

func MakeOtpKeyRegistInput(id string, uuid string, deviceOtp []deviceAuthService.DevicePubKey) input.OtpKeyRegistInput {

	devicePubKeySlice := make([]input.DevicePubKeyInput, 0)

	for i := 0; i < len(deviceOtp); i++ {

		temp := input.DevicePubKeyInput{
			Kind: deviceOtp[i].Kind,
			Key:  deviceOtp[i].Key,
		}

		devicePubKeySlice = append(devicePubKeySlice, temp)
	}

	return input.OtpKeyRegistInput{
		Id:           id,
		Uuid:         uuid,
		DevicePubKey: devicePubKeySlice,
	}
}

func MakeOtpKeyRegistOutput(entity entity.OtpKeyRegistResultEntity) output.OtpKeyRegistOutput {
	return output.OtpKeyRegistOutput{
		OtpRegDate:       entity.OtpRegDate,
		SvChatKeyVersion: entity.SvChatKeyVersion,
		SvNoteKeyVersion: entity.SvNoteKeyVersion,
	}
}
