package entity

type OTPKeyRegistEntity struct {
	Id               string
	Uuid             string
	OtpKeyInfoEntity []OtpKeyInfoEntity
}

func MakeOtpKeyRegistEntity(id string, uuid string, otpKeyInfoEntity []OtpKeyInfoEntity) OTPKeyRegistEntity {

	return OTPKeyRegistEntity{
		Id:               id,
		Uuid:             uuid,
		OtpKeyInfoEntity: otpKeyInfoEntity,
	}
}
