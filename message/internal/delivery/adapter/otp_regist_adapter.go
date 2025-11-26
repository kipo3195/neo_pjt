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

func MakeOtpKeyRegistOutput(chkeyRegDate, noKeyRegDate string) output.OtpKeyregistOutput {
	return output.OtpKeyregistOutput{
		ChkeyRegDate: chkeyRegDate,
		NoKeyRegDate: noKeyRegDate,
	}
}
