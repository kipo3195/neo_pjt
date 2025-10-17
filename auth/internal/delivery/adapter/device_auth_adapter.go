package adapter

import "auth/internal/application/usecase/input"

func MakeDeviceRegistCheckInput(id string, uuid string, challenge string) input.DeviceRegistInput {

	return input.DeviceRegistInput{
		Id:        id,
		Uuid:      uuid,
		Challenge: challenge,
	}
}

// func MakeDeviceAuthTokenRegistInput(id string, uuid string) input.DeviceAuthTokenRegistInput {

// 	return input.DeviceAuthTokenRegistInput{
// 		Id:   id,
// 		Uuid: uuid,
// 	}

// }

func MakeRemoveDeviceChallengeInput(id string, uuid string) input.RemoveDeviceChallengeInput {

	return input.RemoveDeviceChallengeInput{
		Id:   id,
		Uuid: uuid,
	}
}

func MakeUpdateDeviceInfoInput(id string, uuid string, modelName string, version string) input.UpdateDeviceInfoInput {

	return input.UpdateDeviceInfoInput{
		Id:        id,
		Uuid:      uuid,
		ModelName: modelName,
		Version:   version,
	}

}
