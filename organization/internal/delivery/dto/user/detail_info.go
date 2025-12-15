package user

type DetailInfo struct {
	UserHash  string        `json:"userHash"`
	UserName  UserNameDto   `json:"userName"`
	MyOrgCode string        `json:"myOrgCode,omitempty" `
	DeptInfo  []DeptInfoDto `json:"deptInfo"`
}
