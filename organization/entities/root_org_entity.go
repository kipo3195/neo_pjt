package entities

type RootOrg struct {
	OrgInfogs []OrgInfo
}

type OrgInfo struct {
	DeptCode       string `json:"deptCode"`
	ParentDeptCode string `json:"parentDeptCode"`
	KrLang         string `json:"krLang"`
	EnLang         string `json:"enLang"`
	CnLang         string `json:"cnLang"`
	JpLang         string `json:"jpLang"`
	DeptUpdateHash string `json:"deptUpdateHash"`
}
