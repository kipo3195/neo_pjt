package entities

type OrgTreeInfos struct {
	DeptCode       string     `json:"deptCode"`
	ParentDeptCode string     `json:"parentDeptCode"`
	Name           NameEntity `json:"name"`
	//UpdateHash     string         `json:"updateHash"`
	SubDept []OrgTreeInfos `json:"subDept,omitempty"` // 재귀 구조
	Id      string         `json:"id"`
	Kind    string         `json:"kind"` // 사용자, 부서 구분
}
