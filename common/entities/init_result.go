package entities

type InitResult struct {
	IssuedAppToken *IssuedAppToken `json:"issuedAppToken"`
	ConnectInfo    *ConnectInfo    `json:"connectInfo"`
	WorksConfig    *WorksConfig    `json:"worksConfig"`
}
