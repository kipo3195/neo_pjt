package input

type UserRegisterChallengeInput struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
}

func MakeUserRegisterChallengeInput(id string, salt string) UserRegisterChallengeInput {

	return UserRegisterChallengeInput{
		Id:   id,
		Salt: salt,
	}
}
