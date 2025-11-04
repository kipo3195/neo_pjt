package output

type DeptInfoOutput struct {
	DeptOrg      string             `json:"deptOrg"`
	DeptCode     string             `json:"deptCode"`
	DeptName     DeptNameOutput     `json:"deptName"`
	Header       string             `json:"header"`
	Description  string             `json:"description"`
	RoleName     RoleNameOutput     `json:"roleName"`
	PositionName PositionNameOutput `json:"positionName"`
}

type RoleNameOutput struct {
	KoLang string `json:"ko"`
	EnLang string `json:"en"`
	ZhLang string `json:"zh"`
	JpLang string `json:"jp"`
}

type PositionNameOutput struct {
	KoLang string `json:"ko"`
	EnLang string `json:"en"`
	ZhLang string `json:"zh"`
	JpLang string `json:"jp"`
}

type DeptNameOutput struct {
	DefLang string `json:"def"`
	KoLang  string `json:"ko"`
	EnLang  string `json:"en"`
	ZhLang  string `json:"zh"`
	JpLang  string `json:"jp"`
	RuLang  string `json:"ru"`
	ViLang  string `json:"vi"`
}
