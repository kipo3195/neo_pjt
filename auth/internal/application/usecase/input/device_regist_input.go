package input

type DeviceRegistInput struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	ModelName string `json:"modelName"`
	Version   string `json:"version"`
	Challenge string `json:"challenge"`
}

func MakeDeviceRegistInput(id string, uuid string, modelName string, version string, challenge string) DeviceRegistInput {

	return DeviceRegistInput{
		Id:        id,
		Uuid:      uuid,
		ModelName: modelName,
		Version:   version,
		Challenge: challenge,
	}
}
