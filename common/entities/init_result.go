package entities

type InitResult struct {
	AppToken    string       `json:"appToken"`
	ConnectInfo *ConnectInfo `json:"connectInfo"`
	WorksConfig *WorksConfig `json:"worksConfig"`
}
