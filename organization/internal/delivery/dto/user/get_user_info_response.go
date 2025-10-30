package user

type GetUserInfoResponse struct {
	UserHash string        `json:"userHash"`
	Username UsernameDto   `json:"userName"`
	OrgCode  []string      `json:"orgCode"`
	DeptInfo []DeptInfoDto `json:"deptInfo"`
}
