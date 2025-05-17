package entities

type WorksInfo struct {
	WorksCode   string      `json:"worksCode"`
	WorksName   string      `json:"worksName"`
	UseYn       string      `json:"useYn"`
	RegDate     string      `json:"regDate"`
	WorksAuth   WorksAuth   `json:"worksAuth"`
	ConnectInfo ConnectInfo `json:"connectInfo"`
}
