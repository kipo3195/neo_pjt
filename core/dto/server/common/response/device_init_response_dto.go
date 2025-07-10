package common

type DeviceInitResponseBody struct {
	IssuedAppToken *IssuedAppTokenDTO `json:"issuedAppToken"`
	//ConnectInfo *ConnectInfo `json:"connectInfo"`
	// 20250706 connectInfo 즉, 서버 url은 core에서도 갖고 있기때문에 내려주지 않음. 단, common에서 commonInfo를 유지하도록 함 (클라이언트 최초 연결 이후 대응)
	WorksConfig *WorksConfigDTO `json:"worksConfig"`
}

type DeviceInitResponseDTO struct {
	Body *DeviceInitResponseBody
}
