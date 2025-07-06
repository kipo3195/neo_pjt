package common

type DeviceInitResponse struct {
	IssuedAppToken *IssuedAppToken `json:"issuedAppToken"`
	//ConnectInfo *ConnectInfo `json:"connectInfo"`
	// 20250706 connectInfo 즉, 서버 url은 core에서도 갖고 있기때문에 처리하지않음. 단, common에서 commonInfo를 유지하도록 함(클라이언트 최초 연결 이후 대응)
	WorksConfig *WorksConfig `json:"worksConfig"`
}
