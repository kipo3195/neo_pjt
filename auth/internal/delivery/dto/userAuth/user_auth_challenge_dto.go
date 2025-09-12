package userAuth

type UserAuthChallengeRequest struct {
	Id     string `json:"id"`
	Device string `json:"device"`
}
