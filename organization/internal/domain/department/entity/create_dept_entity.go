package entity

type CreateDeptEntity struct {
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

func MakeCreateDeptEntity(deptCode string, deptOrg string, parentDeptCode string, ko string, en string, jp string, ru string, vi string, zh string, header string) CreateDeptEntity {
	return CreateDeptEntity{
		DeptOrg:        deptOrg,
		DeptCode:       deptCode,
		ParentDeptCode: parentDeptCode,
		KoLang:         ko,
		EnLang:         en,
		JpLang:         jp,
		RuLang:         ru,
		ViLang:         vi,
		ZhLang:         zh,
		Header:         header,
	}

}
