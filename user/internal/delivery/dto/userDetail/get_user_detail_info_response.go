package userDetail

type GetUserDetailInfoResponse struct {
	UserHash     string `json:"userHash"`
	UserEmail    string `json:"userEmail"`
	UserPhoneNum string `json:"userPhoneNum"`

	ProfileMsg string `json:"profileMsg"`
}
