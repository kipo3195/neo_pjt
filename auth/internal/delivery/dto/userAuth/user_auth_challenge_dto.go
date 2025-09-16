package userAuth

type UserAuthChallengeRequest struct {
	Id     string `json:"id"`
	Device string `json:"device"`
}

type UserAuthChallengeResponse struct {
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
}
