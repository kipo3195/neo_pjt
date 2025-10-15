package adapter

import "auth/internal/application/usecase/input"

func MakeDeviceRegistStateInput(id string, uuid string) input.DeviceRegistStateInput {

	return input.DeviceRegistStateInput{
		Id:   id,
		Uuid: uuid,
	}

}
