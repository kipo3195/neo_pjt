package user

type GetMyInfoResponse struct {
	UserHash string        `json:"userHash"`
	Username UsernameDto   `json:"userName"`
	OrgCode  []string      `json:"orgCode"`
	DeptInfo []DeptInfoDto `json:"deptInfo"`
}

//userHash
//userName
//orgCode 배열
//profile
//deptInfo
