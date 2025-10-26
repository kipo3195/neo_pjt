package output

type UserAuthChallengeOutput struct {
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
}
