package entities

type OrgTreeInfos struct {
	DeptCode       string         `json:"deptCode"`
	ParentDeptCode string         `json:"parentDeptCode"`
	KrLang         string         `json:"krLang"`
	EnLang         string         `json:"enLang"`
	CnLang         string         `json:"cnLang"`
	JpLang         string         `json:"jpLang"`
	UpdateHash     string         `json:"updateHash"`
	SubDept        []OrgTreeInfos `json:"subDept,omitempty"` // 재귀 구조
	Kind           string         `json:"kind"`              // 사용자, 부서 구분
}
