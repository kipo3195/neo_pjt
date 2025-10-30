package user

type GetUserInfoRequest struct {
	UserHashs []string `json:"userHashs" validate:"required"`
}
