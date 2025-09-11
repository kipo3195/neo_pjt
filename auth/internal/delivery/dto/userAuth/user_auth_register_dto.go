package userAuth

type UserAuthRegisterRequest struct {
	Id       string `json:"id"`
	Salt     string `json:"salt"`
	UserHash string `json:"userHash"`
	AuthHash string `json:"authHash"`
}
