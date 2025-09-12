package entity

type UserAuthEntity struct {
	Id     string
	Fv     string
	Device string
}

func MakeUserAuthEntity(id string, fv string, device string) UserAuthEntity {
	return UserAuthEntity{
		Id:     id,
		Fv:     fv,
		Device: device,
	}
}
