package serviceUser

type RegistServiceUserRequest struct {
	Org          string   `json:"org"`
	UserId       []string `json:"userId"`
	UserIdPrefix string   `json:"userIdPrefix"`
	UserAuth     string   `json:"userAuth"`
	Start        int      `json:"start"`
	End          int      `json:"end"`
}
