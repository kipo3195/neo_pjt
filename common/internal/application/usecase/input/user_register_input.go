package input

type UserRegisterInput struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
	Fv   string `json:"fv"`
}

func MakeUserRegisterInput(id string, salt string, fv string) UserRegisterInput {

	return UserRegisterInput{
		Id:   id,
		Salt: salt,
		Fv:   fv,
	}
}
