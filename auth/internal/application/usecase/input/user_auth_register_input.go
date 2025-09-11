package input

type UserAuthRegisterInput struct {
	Id       string `json:"id"`
	Salt     string `json:"salt"`
	UserHash string `json:"userHash"`
	AuthHash string `json:"authHash"`
}

func MakeUserAuthRegisterInput(id string, salt string, authHash string, userHash string) UserAuthRegisterInput {
	return UserAuthRegisterInput{
		Id:       id,
		Salt:     salt,
		UserHash: userHash,
		AuthHash: authHash,
	}
}
