package device

type DeviceDTO struct {
	Body DeviceRequestBody
}

type DeviceRequestBody struct {
	WorksCode string `json:"worksCode"`
	Uuid      string `json:"uuid"`
	Device    string `json:"device"`
}

type DeviceResultResponse struct {
	WorksInfo      interface{} `json:"worksInfo"`
	IssuedAppToken interface{} `json:"issuedAppToken"`
	SkinInfo       interface{} `json:"skinInfo"`
	WorksConfig    interface{} `json:"worksConfig"`
}
