package input

type UserAuthChallengeInput struct {
	Id     string
	Device string
}

func MakeUserAuthChallengeInput(id string, device string) UserAuthChallengeInput {
	return UserAuthChallengeInput{
		Id:     id,
		Device: device,
	}
}
