package entity

type UserRegisterEntity struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
	Fv   string `json:"fv"`
}

func MakeUserRegisterEntity(id string, salt string, fv string) UserRegisterEntity {

	return UserRegisterEntity{
		Id:   id,
		Salt: salt,
		Fv:   fv,
	}
}
