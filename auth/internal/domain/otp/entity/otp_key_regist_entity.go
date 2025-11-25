package entity

type OtpKeyRegistEntity struct {
	Id    string
	Uuid  string
	ChKey string
	NoKey string
}

func MakeOtpKeyRegistEntity(id string, uuid string, chKey string, noKey string) OtpKeyRegistEntity {
	return OtpKeyRegistEntity{
		Id:    id,
		Uuid:  uuid,
		ChKey: chKey,
		NoKey: noKey,
	}
}
