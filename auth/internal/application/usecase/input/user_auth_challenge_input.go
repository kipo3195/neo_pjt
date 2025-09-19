package input

type UserAuthChallengeInput struct {
	Id   string
	Uuid string
}

func MakeUserAuthChallengeInput(id string, uuid string) UserAuthChallengeInput {
	return UserAuthChallengeInput{
		Id:   id,
		Uuid: uuid,
	}
}
