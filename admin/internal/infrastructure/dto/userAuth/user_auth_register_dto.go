package userAuth

type UserAuthRegisterDto struct {
	UserId   string `json:"userId"`
	UserHash string `json:"userHash"`
	UserAuth string `json:"userAuth"`
	Salt     string `json:"salt"`
}
