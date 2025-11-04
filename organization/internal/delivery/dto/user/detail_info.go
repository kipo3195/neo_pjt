package user

type DetailInfo struct {
	UserHash string        `json:"userHash"`
	UserName UserNameDto   `json:"userName"`
	OrgCode  []string      `json:"orgCode"`
	DeptInfo []DeptInfoDto `json:"deptInfo"`
}
