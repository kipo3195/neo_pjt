package server

type SvCreateDeptRequest struct {
	DeptOrg        string `json:"deptOrg"`
	DeptCode       string `json:"deptCode"`
	ParentDeptCode string `json:"parentDeptCode"`
	KoLang         string `json:"ko"`
	EnLang         string `json:"en"`
	JpLang         string `json:"jp"`
	ZhLang         string `json:"zh"`
	RuLang         string `json:"ru"`
	ViLang         string `json:"vi"`
	Header         string `json:"header"`
}
