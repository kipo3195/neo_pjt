package input

type UserAuthInput struct {
	Id   string
	Fv   string
	Uuid string
}

func MakeUserAuthInput(id string, fv string, uuid string) UserAuthInput {
	return UserAuthInput{
		Id:   id,
		Fv:   fv,
		Uuid: uuid,
	}
}
