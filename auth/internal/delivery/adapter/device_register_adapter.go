package adapter

import "auth/internal/application/usecase/input"

func MakeDeviceRegistInput(id string, uuid string, modelName string, version string, challenge string) input.DeviceRegistInput {

	return input.DeviceRegistInput{
		Id:        id,
		Uuid:      uuid,
		ModelName: modelName,
		Version:   version,
		Challenge: challenge,
	}
}
