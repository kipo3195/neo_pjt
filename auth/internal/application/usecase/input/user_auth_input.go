package input

type UserAuthInput struct {
	Id     string
	Fv     string
	Device string
}

func MakeUserAuthInput(id string, fv string, device string) UserAuthInput {
	return UserAuthInput{
		Id:     id,
		Fv:     fv,
		Device: device,
	}
}
