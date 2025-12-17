package serviceUser

type RegistServiceUserRequest struct {
	Org      string   `json:"org"`
	UserId   []string `json:"userId"`
	UserAuth string   `json:"userAuth"`
}
