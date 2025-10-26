package input

type UserRegisterInput struct {
	Id   string `json:"id"`
	Salt string `json:"salt"`
	Fv   string `json:"fv"`
}
