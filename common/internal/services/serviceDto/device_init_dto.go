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
	WorksInfo      interface{} `json:"worksInfo"`
	IssuedAppToken interface{} `json:"issuedAppToken"`
	SkinInfo       interface{} `json:"skinInfo"`
	WorksConfig    interface{} `json:"worksConfig"`
}
