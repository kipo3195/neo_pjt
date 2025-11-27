package adapter

import (
	"message/internal/application/usecase/input"
	"message/internal/application/usecase/output"
)

func MakeOtpKeyRegistInput(id, uuid, chKey, noKey string) input.OtpKeyRegistInput {
	return input.OtpKeyRegistInput{
		Id:    id,
		Uuid:  uuid,
		ChKey: chKey,
		NoKey: noKey,
	}
}

func MakeOtpKeyRegistOutput(otpRegDate string, serverKeyVersion string) output.OtpKeyregistOutput {
	return output.OtpKeyregistOutput{
		OtpRegDate:   otpRegDate,
		SvKeyVersion: serverKeyVersion,
	}
}
