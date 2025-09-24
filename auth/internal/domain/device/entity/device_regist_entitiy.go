package entity

type DeviceRegistEntity struct {
	Id        string `json:"id"`
	Uuid      string `json:"uuid"`
	Version   string `json:"version"`
	ModelName string `json:"modelName"`
	Challenge string `json:"challenge"`
}

func MakeDeviceRegistEntity(id string, uuid string, modelName string, version string, challenge string) DeviceRegistEntity {
	return DeviceRegistEntity{
		Id:        id,
		Uuid:      uuid,
		ModelName: modelName,
		Version:   version,
		Challenge: challenge,
	}
}
