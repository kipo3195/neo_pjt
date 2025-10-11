package entity

type RegisterDeptEntity struct {
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

func MakeRegisterDeptEntity(deptOrg string,
	deptCode string, parentDeptCode string, koLang string, enLang string, jpLang string, zhLang string, ruLang string, viLang string, header string) RegisterDeptEntity {
	return RegisterDeptEntity{
		DeptCode:       deptCode,
		DeptOrg:        deptOrg,
		ParentDeptCode: parentDeptCode,
		KoLang:         koLang,
		EnLang:         enLang,
		JpLang:         jpLang,
		ZhLang:         zhLang,
		RuLang:         ruLang,
		ViLang:         viLang,
		Header:         header,
	}
}
