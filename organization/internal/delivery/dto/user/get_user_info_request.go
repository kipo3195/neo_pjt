package user

type GetUserInfoRequest struct {
	UserIds []string `json:"userIds" validate:"required"`
}
