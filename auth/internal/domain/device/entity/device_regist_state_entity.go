package entity

type DeviceRegistStateEntity struct {
	Id   string `json:"id"`
	Uuid string `json:"uuid"`
}

func MakeDeviceRegistStateEntity(id string, uuid string) DeviceRegistStateEntity {
	return DeviceRegistStateEntity{
		Id:   id,
		Uuid: uuid,
	}
}
