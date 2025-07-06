package entities

type WorksInfo struct {
	WorksCode    string               `json:"worksCode"`
	WorksName    string               `json:"worksName"`
	UseYn        string               `json:"useYn"`
	RegDate      string               `json:"regDate"`
	WorksAuth    WorksAuth            `json:"worksAuth"`
	ConnectInfo  ConnectInfo          `json:"connectInfo"`
	TimeZone     string               `json:"timeZone"`   // 기본 시간대
	Language     string               `json:"language"`   // 기본 언어
	SkinHash     string               `json:"skinHash"`   // 현재 앱 스킨의 버전
	ConfigHash   string               `json:"configHash"` // 설정의 버전, 클라이언트 별로 구분됨
	SkinFileInfo []SkinFileInfoEntity `json:"skin"`
}
