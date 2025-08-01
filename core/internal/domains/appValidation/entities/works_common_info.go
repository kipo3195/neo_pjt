package entities

// 일반 설정
type WorksCommonInfo struct {
	ServerUrl string `json:"serverUrl"`
	WorksCode string `json:"worksCode"`
	WorksName string `json:"worksName"`
	UseYn     string `json:"useYn"`
	RegDate   string `json:"regDate"`
}
