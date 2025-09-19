package entity

type UserAuthEntity struct {
	Id   string
	Fv   string
	Uuid string
}

func MakeUserAuthEntity(id string, fv string, uuid string) UserAuthEntity {
	return UserAuthEntity{
		Id:   id,
		Fv:   fv,
		Uuid: uuid,
	}
}
