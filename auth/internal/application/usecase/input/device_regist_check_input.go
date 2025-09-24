package input

type DeviceRegistCheckInput struct {
	Id   string `json:"id"`
	Uuid string `json:"uuid"`
}

func MakeDeviceRegistCheckInput(id string, uuid string) DeviceRegistCheckInput {

	return DeviceRegistCheckInput{
		Id:   id,
		Uuid: uuid,
	}
}
