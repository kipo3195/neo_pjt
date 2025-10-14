package input

type AccessTokenReIssueInput struct {
	AppToken     string
	Uuid         string
	RefreshToken string
}

func MakeAccessTokenReIssueInput(appToken string, uuid string, refreshToken string) AccessTokenReIssueInput {
	return AccessTokenReIssueInput{
		AppToken:     appToken,
		Uuid:         uuid,
		RefreshToken: refreshToken,
	}
}
