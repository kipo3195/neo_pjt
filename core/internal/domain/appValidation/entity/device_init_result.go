package entity

type DeviceInitResult struct {
	IssuedAppToken *IssuedAppToken `json:"issuedAppToken"`
	WorksConfig    interface{}     `json:"worksConfig"`
}
