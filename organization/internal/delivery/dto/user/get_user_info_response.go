package user

type GetUserInfoResponse struct {
	UserHash   string        `json:"userHash"`
	UserDetail UserDetail    `json:"userDetail"`
	Username   UsernameDto   `json:"userName"`
	OrgCode    []string      `json:"orgCode"`
	DeptInfo   []DeptInfoDto `json:"deptInfo"`
}
