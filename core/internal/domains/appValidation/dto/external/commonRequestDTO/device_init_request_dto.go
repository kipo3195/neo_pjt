package commonRequestDTO

type DeviceInitRequestHeader struct {
	ServerToken string `json:"serverKey"`
}

type DeviceInitRequestBody struct {
	Uuid      string `json:"uuid"`
	WorksCode string `json:"worksCode"`
	Device    string `json:"device"`
}

type DeviceInitRequestDTO struct {
	Header DeviceInitRequestHeader
	Body   DeviceInitRequestBody
}
