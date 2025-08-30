package user

type GetMyInfoResponseDTO struct {
	Body GetMyInfoResponseBody
}

type GetMyInfoResponseBody struct {
	UserHash     string        `json:"userHash"`
	UserPhoneNum string        `json:"userPhoneNum"`
	Username     UsernameDto   `json:"userName"`
	OrgCodes     []string      `json:"orgCodes"`
	ProfileUrl   string        `json:"profileUrl"`
	ProfileMsg   string        `json:"profileMsg"`
	DeptInfo     []DeptInfoDto `json:"deptInfo"`
}
