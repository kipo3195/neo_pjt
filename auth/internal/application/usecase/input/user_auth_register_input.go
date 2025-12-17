package input

type UserAuthRegisterInput struct {
	UserId   string `json:"userId"`
	Salt     string `json:"salt"`
	UserHash string `json:"userHash"`
	UserAuth string `json:"userAuth"`
}
