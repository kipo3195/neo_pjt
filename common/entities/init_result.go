package entities

type InitResult struct {
	AppToken      string `json:"appToken"`
	ConnectInfo   string `json:"connectInfo"`
	TimeZone      string `json:"timeZone"`      // 기본 시간대
	Language      string `json:"language"`      // 기본 언어
	SkinVersion   string `json:"skinVersion"`   // 현재 앱 스킨의 버전 클라이언트 별로 구분됨.
	ConfigVersion string `json:"configVersion"` // 설정의 버전, 클라이언트 공통

}
