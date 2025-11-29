package entity

type OTPKeyRegistEntity struct {
	Id           string
	Uuid         string
	ChKey        string
	NoKey        string
	ChatOtpKey   string
	NoteOtpKey   string
	OtpRegDate   string
	SvKeyVersion string
}

func MakeOtpKeyRegistEntity(id, uuid, chKey, noKey string) OTPKeyRegistEntity {

	return OTPKeyRegistEntity{
		Id:    id,
		Uuid:  uuid,
		ChKey: chKey,
		NoKey: noKey,
	}
}
