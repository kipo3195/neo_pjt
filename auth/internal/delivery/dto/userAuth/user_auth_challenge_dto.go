package userAuth

type UserAuthChallengeRequest struct {
	Id   string `json:"id" validate:"required"`
	Uuid string `json:"uuid" validate:"required"`
}

type UserAuthChallengeResponse struct {
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
}
