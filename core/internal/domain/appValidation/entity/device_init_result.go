package entity

type DeviceInitResult struct {
	WorksCommonInfo WorksCommonInfo
	IssuedAppToken  *IssuedAppToken `json:"issuedAppToken"`
	WorksConfig     interface{}     `json:"worksConfig"`
}
