package user

type UserRegisterChallengeRequest struct {
	Id   string `json:"id" validate:"required"`
	Salt string `json:"string" validate:"required"`
}

type UserRegisterChallengeResponse struct {
	Challenge string `json:"challenge"`
}
