package user

type MyDetailInfo struct {
	UserHash  string        `json:"userHash"`
	UserName  UserNameDto   `json:"userName"`
	MyOrgCode string        `json:"myOrgCode,omitempty" `
	DeptInfo  []DeptInfoDto `json:"deptInfo"`
}
