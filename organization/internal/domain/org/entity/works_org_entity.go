package entity

type WorksOrg struct {
	Org            string `json:"org"` // 추가
	DeptCode       string `json:"deptCode"`
	ParentDeptCode string `json:"parentDeptCode"`
	KoLang         string `json:"koLang"`
	EnLang         string `json:"enLang"`
	ZhLang         string `json:"zhLang"`
	JpLang         string `json:"jpLang"`
	RuLang         string `json:"ruLang"`
	ViLang         string `json:"viLang"`
	UpdateHash     string `json:"updateHash"`
	Kind           string `json:"kind"` // 추가
	UserHash       string `json:"userHash"`
	Header         string `json:"header"`
}
