package output

type UserAuthChallengeOutput struct {
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
}

func MakeUserAuthChallengeOutput(c string, s string) UserAuthChallengeOutput {
	return UserAuthChallengeOutput{
		Challenge: c,
		Salt:      s,
	}
}
