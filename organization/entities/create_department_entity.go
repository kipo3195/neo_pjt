package entities

type CreateDepartmentEntity struct {
	DeptOrg        string `json:"deptOrg"`
	DeptCode       string `json:"deptCode"`
	ParentDeptCode string `json:"parentDeptCode"`
	KrLang         string `json:"krLang"`
	EnLang         string `json:"enLang"`
	JpLang         string `json:"jpLang"`
	CnLang         string `json:"cnLang"`
}
