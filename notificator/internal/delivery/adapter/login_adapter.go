package adapter

import "notificator/internal/application/usecase/input"

func MakeLoginInput(uuid string, deviceType string) input.LoginInput {

	return input.LoginInput{
		Uuid:       uuid,
		DeviceType: deviceType,
	}
}
