package adapter

import "auth/internal/application/usecase/input"

func MakeDeviceRegistCheckInput(id string, uuid string) input.DeviceRegistInput {

	return input.DeviceRegistInput{
		Id:   id,
		Uuid: uuid,
	}
}

// func MakeDeviceAuthTokenRegistInput(id string, uuid string) input.DeviceAuthTokenRegistInput {

// 	return input.DeviceAuthTokenRegistInput{
// 		Id:   id,
// 		Uuid: uuid,
// 	}

// }
