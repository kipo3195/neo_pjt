package user

type DetailInfo struct {
	UserHash string        `json:"userHash"`
	UserName UserNameDto   `json:"userName"`
	DeptInfo []DeptInfoDto `json:"deptInfo"`
}
