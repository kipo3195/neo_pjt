package entity

// 20251214
// batch에서 조회 된 결과를 전송하는 것과 동일한 구조 (view table 조회 결과를 담는다. )
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
	UserId         string `json:"userId"`
	Header         string `json:"header"`
	Description    string `json:"description"`
}
