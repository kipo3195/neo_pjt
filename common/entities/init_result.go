package entities

type InitResult struct {
	IssuedAppToken *IssuedAppToken `json:"IssuedAppToken"`
	ConnectInfo    *ConnectInfo    `json:"connectInfo"`
	WorksConfig    *WorksConfig    `json:"worksConfig"`
}
