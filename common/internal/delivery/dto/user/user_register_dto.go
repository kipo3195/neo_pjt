package user

type UserRegisterRequest struct {
	Id   string `json:"id" validate:"required"`
	Salt string `json:"salt" validate:"required"`
	Fv   string `json:"fv" validate:"required"`
}
