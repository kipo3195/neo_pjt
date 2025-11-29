package entity

type OtpKeyRegistEntity struct {
	Id                 string
	Uuid               string
	DevicePubKeyEntity []DevicePubKeyEntity
}

func MakeOtpKeyRegistEntity(id string, uuid string, devicePubKeyEntity []DevicePubKeyEntity) OtpKeyRegistEntity {
	return OtpKeyRegistEntity{
		Id:                 id,
		Uuid:               uuid,
		DevicePubKeyEntity: devicePubKeyEntity,
	}
}
