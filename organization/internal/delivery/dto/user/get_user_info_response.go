package user

type GetUserInfoResponse struct {
	UserHash   string        `json:"userHash"`
	UserDetail UserDetail    `json:"userDetail"`
	Username   UsernameDto   `json:"userName"`
	Profile    UserProfile   `json:"profile"`
	OrgCode    []string      `json:"orgCode"`
	DeptInfo   []DeptInfoDto `json:"deptInfo"`
}
