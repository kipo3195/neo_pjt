package serviceDto

type DeviceInitDTO struct {
	Body DeviceInitRequestBody
}

type DeviceInitRequestBody struct {
	WorksCode string `json:"worksCode"`
	Uuid      string `json:"uuid"`
	Device    string `json:"device"`
}

type DeviceInitResultResponse struct {
	IssuedAppToken interface{} `json:"issuedAppToken"`
	ConnectInfo    interface{} `json:"connectInfo"`
	WorksConfig    interface{} `json:"worksConfig"`
}
