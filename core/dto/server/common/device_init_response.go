package common

type DeviceInitResponse struct {
	AppToken    string       `json:"appToken"`
	ConnectInfo *ConnectInfo `json:"connectInfo"`
	WorksConfig *WorksConfig `json:"worksConfig"`
}
