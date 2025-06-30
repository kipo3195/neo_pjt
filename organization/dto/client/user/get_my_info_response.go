package client

type GetMyInfoResponse struct {
	UserHash     string      `json:"userHash"`
	UserPhoneNum string      `json:"userPhoneNum"`
	Username     UsernameDto `json:"userName"`
	OrgCodes     []string    `json:"orgCodes"`
	ProfileUrl   string      `json:"profileUrl"`
	ProfileMsg   string      `json:"profileMsg"`
}
