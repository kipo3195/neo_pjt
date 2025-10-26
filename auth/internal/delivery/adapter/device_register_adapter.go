package adapter

import (
	"auth/internal/application/usecase/input"
	"auth/internal/application/usecase/output"
)

func MakeDeviceRegistInput(id string, uuid string, modelName string, version string, challenge string) input.DeviceRegistInput {

	return input.DeviceRegistInput{
		Id:        id,
		Uuid:      uuid,
		ModelName: modelName,
		Version:   version,
		Challenge: challenge,
	}
}

func MakeDeviceRegistOutput(at string, rt string, rtExp string) output.DeviceRegistOutput {
	return output.DeviceRegistOutput{
		AccessToken:     at,
		RefreshToken:    rt,
		RefreshTokenExp: rtExp,
	}
}
