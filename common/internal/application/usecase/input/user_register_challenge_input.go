package input

type UserRegisterChallengeInput struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
}
