package user

type DeptInfoDto struct {
	DeptOrg     string `json:"deptOrg"`
	DeptCode    string `json:"deptCode"`
	DefLang     string `json:"def"`
	KoLang      string `json:"ko"`
	EnLang      string `json:"en"`
	ZhLang      string `json:"zh"`
	JpLang      string `json:"jp"`
	RuLang      string `json:"ru"`
	ViLang      string `json:"vi"`
	Header      string `json:"header"`
	Description string `json:"description"`
}
