package mapper

import (
	"message/internal/adapter/http/dto/otp"
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
)

func MakeOtpKeyRegistInput(id string, uuid string, devicePubKey []otp.DevicePubKey) input.OtpKeyRegistInput {

	devicePubKeySlice := make([]input.DevicePubKeyInput, 0)

	for i := 0; i < len(devicePubKey); i++ {
		in := input.DevicePubKeyInput{
			Kind: devicePubKey[i].Kind,
			Key:  devicePubKey[i].Key,
		}

		devicePubKeySlice = append(devicePubKeySlice, in)
	}

	return input.OtpKeyRegistInput{
		Id:           id,
		Uuid:         uuid,
		DevicePubKey: devicePubKeySlice,
	}
}

func MakeOtpKeyRegistOutput(otpRegDate string, svChatKeyVersion string, svNoteKeyVersion string) output.OtpKeyregistOutput {
	return output.OtpKeyregistOutput{
		OtpRegDate:       otpRegDate,
		SvChatKeyVersion: svChatKeyVersion,
		SvNoteKeyVersion: svNoteKeyVersion,
	}
}
